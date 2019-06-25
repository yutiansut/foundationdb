#include "fdbclient/NativeAPI.actor.h"
#include "fdbserver/TesterInterface.actor.h"
#include "fdbserver/workloads/workloads.actor.h"
#include "fdbserver/workloads/BulkSetup.actor.h"
#include "fdbclient/ReadYourWrites.h"
#include "flow/actorcompiler.h"


enum {OP_GETREADVERSION, OP_GET, OP_GETRANGE, OP_SGET, OP_SGETRANGE, OP_UPDATE, OP_INSERT, OP_INSERTRANGE, OP_CLEAR, OP_SETCLEAR, OP_CLEARRANGE, OP_SETCLEARRANGE, OP_COMMIT, MAX_OP};
enum {OP_COUNT, OP_RANGE};
constexpr int MAXKEYVALUESIZE = 1000;
constexpr int RANGELIMIT = 10000;
struct MakoWorkload : TestWorkload {
	uint64_t rowCount, seqNumLen, sampleSize, actorCountPerClient, keyBytes, maxValueBytes, minValueBytes;
	double testDuration, loadTime, warmingDelay, maxInsertRate, transactionsPerSecond, allowedLatency, periodicLoggingInterval;
	bool enableLogging, commitGet, populateData, runBenchmark, preserveData;
	PerfIntCounter xacts, retries, conflicts, commits, totalOps;
	std::vector<PerfIntCounter> opCounters;
	std::vector<uint64_t> insertionCountsToMeasure;
	std::vector<std::pair<uint64_t, double>> ratesAtKeyCounts;
	std::string operationsSpec;
	//store operations to execute
	int operations[MAX_OP][2];
	// used for periodically tracing
	std::vector<PerfMetric> periodicMetrics;
	// store latency of each operation with sampling
	std::vector<ContinuousSample<double>> opLatencies;
	// prefix of keys populated, e.g. 'mako00000xxxxxxx'
	const std::string KEYPREFIX = "mako";
	const int KEYPREFIXLEN = KEYPREFIX.size();
	const std::array<std::string, MAX_OP> opNames = {"GRV", "GET", "GETRANGE", "SGET", "SGETRANGE", "UPDATE", "INSERT", "INSERTRANGE", "CLEAR", "SETCLEAR", "CLEARRANGE", "SETCLEARRANGE", "COMMIT"};
	MakoWorkload(WorkloadContext const& wcx)
	: TestWorkload(wcx),
	xacts("Transactions"), retries("Retries"), conflicts("Conflicts"), commits("Commits"), totalOps("Operations"),
	loadTime(0.0)
	{
		// init parameters from test file
		// Number of rows populated
		rowCount = getOption(options, LiteralStringRef("rows"), 10000);
		// Test duration in seconds
		testDuration = getOption(options, LiteralStringRef("testDuration"), 30.0);
		warmingDelay = getOption(options, LiteralStringRef("warmingDelay"), 0.0);
		maxInsertRate = getOption(options, LiteralStringRef("maxInsertRate"), 1e12);
		// Flag to control whether to populate data into database
		populateData = getOption(options, LiteralStringRef("populateData"), true);
		// Flag to control whether to run benchmark
		runBenchmark = getOption(options, LiteralStringRef("runBenchmark"), true);
		// Flag to control whether to clean data in the database
		preserveData = getOption(options, LiteralStringRef("preserveData"), true);
		// If true, force commit for read-only transactions
		commitGet = getOption(options, LiteralStringRef("commitGet"), false);
		// Target total transaction-per-second (TPS) of all clients
		transactionsPerSecond = getOption(options, LiteralStringRef("transactionsPerSecond"), 100000.0) / clientCount;
		actorCountPerClient = getOption(options, LiteralStringRef("actorCountPerClient"), 16);
		// Sampling rate (1 sample / <sampleSize> ops) for latency stats
		sampleSize = getOption(options, LiteralStringRef("sampleSize"), rowCount / 100);
		// If true, record latency metrics per periodicLoggingInterval; For details, see tracePeriodically()
		enableLogging = getOption(options, LiteralStringRef("enableLogging"), false);
		periodicLoggingInterval = getOption( options, LiteralStringRef("periodicLoggingInterval"), 5.0 );
		// Specified length of keys and length range of values
		keyBytes = std::max( getOption( options, LiteralStringRef("keyBytes"), 16 ), 16);
		maxValueBytes = getOption( options, LiteralStringRef("valueBytes"), 16 );
		minValueBytes = getOption( options, LiteralStringRef("minValueBytes"), maxValueBytes);
		ASSERT(minValueBytes <= maxValueBytes);
		// The inserted key is formatted as: fixed prefix('mako') + sequential number + padding('x')
		// assume we want to insert 10000 rows with keyBytes set to 16, 
		// then the key goes from 'mako00000xxxxxxx' to 'mako09999xxxxxxx'
		seqNumLen = digits(rowCount);
		// check keyBytes, maxValueBytes is valid
		ASSERT(seqNumLen + KEYPREFIXLEN <= keyBytes);
		ASSERT(keyBytes <= MAXKEYVALUESIZE);
		ASSERT(maxValueBytes <= MAXKEYVALUESIZE);
		// user input: a sequence of operations to be executed; e.g. "g10i5" means to do GET 10 times and Insert 5 times
		// One operation type is defined as "<Type><Count>" or "<Type><Count>:<Range>".
		// When Count is omitted, it's equivalent to setting it to 1.  (e.g. "g" is equivalent to "g1")
		// Multiple operation types can be concatenated.  (e.g. "g9u1" = 9 GETs and 1 update)
		// For RANGE operations, "Range" needs to be specified in addition to "Count".
		// Below are all allowed inputs:
		 	// g – GET
			// gr – GET RANGE
			// sg – Snapshot GET
			// sgr – Snapshot GET RANGE
			// u – Update (= GET followed by SET)
			// i – Insert (= SET with a new key)
			// ir – Insert Range (Sequential)
			// c – CLEAR
			// sc – SET & CLEAR
			// cr – CLEAR RANGE
			// scr – SET & CLEAR RANGE
			// grv – GetReadVersion()
		// Every transaction is committed unless it contains only GET / GET RANGE operations.
		operationsSpec = getOption(options, LiteralStringRef("operations"), LiteralStringRef("g100")).contents().toString();
		//  parse the sequence and extract operations to be executed
		parseOperationsSpec();
		for (int i = 0; i < MAX_OP; ++i) {
			// initilize per-operation latency record
			opLatencies.push_back(ContinuousSample<double>(rowCount / sampleSize));
			// initialize per-operation counter
			opCounters.push_back(PerfIntCounter(opNames[i]));
		}
	}

	std::string description() override {
		// Mako is a simple workload to measure the performance of FDB.
		// The primary purpose of this benchmark is to generate consistent performance results
		return "Mako";
	}

	Future<Void> setup(Database const& cx) override {
		// use all the clients to populate data
		if (populateData)
			return _setup(cx, this);
		return Void();
	}

	Future<Void> start(Database const& cx) override {
		return _start(cx, this);
	}

	Future<bool> check(Database const& cx) override {
		return true;
	}

	void getMetrics(std::vector<PerfMetric>& m) override {
		// metrics of population process
		if (populateData){
			m.push_back( PerfMetric( "Mean load time (seconds)", loadTime, true ) );
			// The importing rate of keys, controlled by parameter "insertionCountsToMeasure"
			auto ratesItr = ratesAtKeyCounts.begin();
			for(; ratesItr != ratesAtKeyCounts.end(); ratesItr++){
				m.push_back(PerfMetric(format("%ld keys imported bytes/sec", ratesItr->first), ratesItr->second, false));
			}
		}
		// benchmark
		if (runBenchmark){
			m.push_back(PerfMetric("Measured Duration", testDuration, true));
			m.push_back(xacts.getMetric());
			m.push_back(PerfMetric("Transactions/sec", xacts.getValue() / testDuration, true));
			m.push_back(totalOps.getMetric());
			m.push_back(PerfMetric("Operations/sec", totalOps.getValue() / testDuration, true));
			m.push_back(conflicts.getMetric());
			m.push_back(PerfMetric("Conflicts/sec", conflicts.getValue() / testDuration, true));
			m.push_back(retries.getMetric());

			// count of each operation
			for (int i = 0; i < MAX_OP; ++i){
				m.push_back(opCounters[i].getMetric());
			}

			// Meaningful Latency metrics
			const int opExecutedAtOnce[] = {OP_GETREADVERSION, OP_GET, OP_GETRANGE, OP_SGET, OP_SGETRANGE, OP_COMMIT};
			for (const int& op : opExecutedAtOnce){
				m.push_back(PerfMetric("Mean " + opNames[op] +" Latency (ms)", 1000 * opLatencies[op].mean(), true));
				m.push_back(PerfMetric("Max " + opNames[op] + " Latency (ms, averaged)", 1000 * opLatencies[op].max(), true));
				m.push_back(PerfMetric("Min " + opNames[op] + " Latency (ms, averaged)", 1000 * opLatencies[op].min(), true));
			}

			//insert logging metrics if exists
			m.insert(m.end(), periodicMetrics.begin(), periodicMetrics.end());
		}
	}
	static std::string randStr(int len) {
		std::string result(len, '.');
		for (int i = 0; i < len; ++i) {
			result[i] = deterministicRandom()->randomAlphaNumeric();
		}
		return result;
	}

	static void randStr(char *str, int len){
		for (int i = 0; i < len; ++i) {
			str[i] = deterministicRandom()->randomAlphaNumeric();
		}
	}

	Value randomValue() {
		const int length = deterministicRandom()->randomInt(minValueBytes, maxValueBytes + 1);
		std::string valueString = randStr(length);
		return StringRef(reinterpret_cast<const uint8_t*>(valueString.c_str()), length);
	}

	Key keyForIndex(uint64_t ind) {
		Key result = makeString(keyBytes);
		char* data = reinterpret_cast<char*>(mutateString(result));
		format((KEYPREFIX + "%0*d").c_str(), seqNumLen, ind).copy(data, KEYPREFIXLEN + seqNumLen);
		for (int i = KEYPREFIXLEN + seqNumLen; i < keyBytes; ++i)
			data[i] = 'x';
		return result;
	}

	/* number of digits */
	static uint64_t digits(uint64_t num) {
		uint64_t digits = 0;
		while (num > 0) {
			num /= 10;
			digits++;
		}
		return digits;
	}
	Standalone<KeyValueRef> operator()(uint64_t n) {
		return KeyValueRef(keyForIndex(n), randomValue());
	}

	ACTOR static Future<Void> tracePeriodically( MakoWorkload *self){
		state double start = now();
		state double elapsed = 0.0;
		state int64_t last_ops = 0;
		state int64_t last_xacts = 0;

		loop {
			elapsed += self->periodicLoggingInterval;
			wait( delayUntil(start + elapsed));
			TraceEvent((self->description() + "_CommitLatency").c_str()).detail("Mean", self->opLatencies[OP_COMMIT].mean()).detail("Median", self->opLatencies[OP_COMMIT].median()).detail("Percentile5", self->opLatencies[OP_COMMIT].percentile(.05)).detail("Percentile95", self->opLatencies[OP_COMMIT].percentile(.95)).detail("Count", self->opCounters[OP_COMMIT].getValue()).detail("Elapsed", elapsed);
			TraceEvent((self->description() + "_GRVLatency").c_str()).detail("Mean", self->opLatencies[OP_GETREADVERSION].mean()).detail("Median", self->opLatencies[OP_GETREADVERSION].median()).detail("Percentile5", self->opLatencies[OP_GETREADVERSION].percentile(.05)).detail("Percentile95", self->opLatencies[OP_GETREADVERSION].percentile(.95)).detail("Count", self->opCounters[OP_GETREADVERSION].getValue());
			
			std::string ts = format("T=%04.0fs: ", elapsed);
			self->periodicMetrics.push_back(PerfMetric(ts + "Transactions/sec", (self->xacts.getValue() - last_xacts) / self->periodicLoggingInterval, false));
			self->periodicMetrics.push_back(PerfMetric(ts + "Operations/sec", (self->totalOps.getValue() - last_ops) / self->periodicLoggingInterval, false));

			last_xacts = self->xacts.getValue();
			last_ops = self->totalOps.getValue();
		}
	}
	ACTOR Future<Void> _setup(Database cx, MakoWorkload* self) {

		state Promise<double> loadTime;
		state Promise<std::vector<std::pair<uint64_t, double>>> ratesAtKeyCounts;

		wait(bulkSetup(cx, self, self->rowCount, loadTime, self->insertionCountsToMeasure.empty(), self->warmingDelay,
		               self->maxInsertRate, self->insertionCountsToMeasure, ratesAtKeyCounts));

		// This is the setup time 
		self->loadTime = loadTime.getFuture().get();
		// This is the rates of importing keys
		self->ratesAtKeyCounts = ratesAtKeyCounts.getFuture().get();

		return Void();
	}

	ACTOR Future<Void> _start(Database cx, MakoWorkload* self) {
		// TODO: Do I need to read data to warm the cache of the keySystem like ReadWrite.actor.cpp (line 465)?
		if (self->runBenchmark) {
			wait(self->_runBenchmark(cx, self));
		}
		if (!self->preserveData && self->clientId == 0){
			wait(self->cleanup(cx, self));
		}
		return Void();
	}

	ACTOR Future<Void> _runBenchmark(Database cx, MakoWorkload* self){
		std::vector<Future<Void>> clients;
		for (int c = 0; c < self->actorCountPerClient; ++c) {
			clients.push_back(self->makoClient(cx, self, self->actorCountPerClient / self->transactionsPerSecond, c));
		}

		if (self->enableLogging)
			clients.push_back(tracePeriodically(self));

		wait( timeout( waitForAll( clients ), self->testDuration, Void() ) );
		return Void();
	}

	ACTOR Future<Void> makoClient(Database cx, MakoWorkload* self, double delay, int actorIndex) {

		state Key rkey, rkey2;
		state Value rval;
		state ReadYourWritesTransaction tr(cx);
		state bool doCommit;
		state int i, count;
		state uint64_t range, indBegin, indEnd, rangeLen;
		state double lastTime = now();
		state double commitStart;
		state KeyRangeRef rkeyRangeRef;
		state std::vector<int> perOpCount(MAX_OP, 0);

		TraceEvent("ClientStarting").detail("ActorIndex", actorIndex).detail("ClientIndex", self->clientId).detail("NumActors", self->actorCountPerClient);

		loop {
			// used for throttling
			wait(poisson(&lastTime, delay));
			try{
				// user-defined value: whether commit read-only ops or not; default is false
				doCommit = self->commitGet;
				for (i = 0; i < MAX_OP; ++i) {
					if (i == OP_COMMIT) 
						continue;
					for (count = 0; count < self->operations[i][0]; ++count) {
						range = std::min(RANGELIMIT, self->operations[i][1]);
						rangeLen = digits(range);
						// generate random key-val pair for operation
						indBegin = self->getRandomKey(self->rowCount);
						rkey = self->keyForIndex(indBegin);
						rval = self->randomValue();
						indEnd = std::min(indBegin + range, self->rowCount);
						rkey2 = self->keyForIndex(indEnd);
						// KeyRangeRef(min, maxPlusOne)
						rkeyRangeRef = KeyRangeRef(rkey, rkey2);

						if (i == OP_GETREADVERSION){
							wait(logLatency(tr.getReadVersion(), &self->opLatencies[i]));
						}
						else if (i == OP_GET){
							wait(logLatency(tr.get(rkey, false), &self->opLatencies[i]));
						} else if (i == OP_GETRANGE){
							wait(logLatency(tr.getRange(rkeyRangeRef, RANGELIMIT, false), &self->opLatencies[i]));
						}
						else if (i == OP_SGET){
							wait(logLatency(tr.get(rkey, true), &self->opLatencies[i]));
						} else if (i == OP_SGETRANGE){
							//do snapshot get range here
							wait(logLatency(tr.getRange(rkeyRangeRef, RANGELIMIT, true), &self->opLatencies[i]));
						} else if (i == OP_UPDATE){
							wait(logLatency(tr.get(rkey, false), &self->opLatencies[OP_GET]));
							tr.set(rkey, rval);
							doCommit = true;
						} else if (i == OP_INSERT){
							// generate an (almost) unique key here, it starts with 'mako' and then comes with randomly generated characters
							randStr(reinterpret_cast<char*>(mutateString(rkey)) + self->KEYPREFIXLEN, self->keyBytes-self->KEYPREFIXLEN);
							tr.set(rkey, rval);
							doCommit = true;
						} else if (i == OP_INSERTRANGE){
							char *rkeyPtr = reinterpret_cast<char*>(mutateString(rkey));
							randStr(rkeyPtr + self->KEYPREFIXLEN, self->keyBytes-self->KEYPREFIXLEN);
							for (int range_i = 0; range_i < range; ++range_i){
								format("%0.*d", rangeLen, range_i).copy(rkeyPtr + self->keyBytes - rangeLen, rangeLen);
								tr.set(rkey, self->randomValue());
							}
							doCommit = true;
						} else if (i == OP_CLEAR){
							tr.clear(rkey);
							doCommit = true;
						} else if(i == OP_SETCLEAR){
							randStr(reinterpret_cast<char*>(mutateString(rkey)) + self->KEYPREFIXLEN, self->keyBytes-self->KEYPREFIXLEN);
							tr.set(rkey, rval);
							// commit the change and update metrics
							commitStart = now();
							wait(tr.commit());
							self->opLatencies[OP_COMMIT].addSample(now() - commitStart);
							++perOpCount[OP_COMMIT];
							tr.reset();
							tr.clear(rkey);
							doCommit = true;
						} else if (i == OP_CLEARRANGE){
							tr.clear(rkeyRangeRef);
							doCommit = true;
						} else if (i == OP_SETCLEARRANGE){
							char *rkeyPtr = reinterpret_cast<char*>(mutateString(rkey));
							randStr(rkeyPtr + self->KEYPREFIXLEN, self->keyBytes-self->KEYPREFIXLEN);
							state std::string scr_start_key;
							state std::string scr_end_key;
							for (int range_i = 0; range_i < range; ++range_i){
								format("%0.*d", rangeLen, range_i).copy(rkeyPtr + self->keyBytes - rangeLen, rangeLen);
								tr.set(rkey, self->randomValue());
								if (range_i == 0)
									scr_start_key = rkey.toString();
							}
							scr_end_key = rkey.toString();
							commitStart = now();
							wait(tr.commit());
							self->opLatencies[OP_COMMIT].addSample(now() - commitStart);
							++perOpCount[OP_COMMIT];
							tr.reset();
							tr.clear(KeyRangeRef(StringRef(scr_start_key), StringRef(scr_end_key)));
							doCommit = true;
						}
						++perOpCount[i];
					}
				}

				if (doCommit) {
					commitStart = now();
					wait(tr.commit());
					self->opLatencies[OP_COMMIT].addSample(now() - commitStart);
					++perOpCount[OP_COMMIT];
				}
				// successfully finish the transaction, update metrics
				++self->xacts;
				for (int op = 0; op < MAX_OP; ++op){
					self->opCounters[op] += perOpCount[op];
					self->totalOps += perOpCount[op];
				}
			} catch (Error& e) {
				TraceEvent("FailedToExecOperations").error(e);
				if (e.code() == error_code_operation_cancelled)
					throw;
				else if (e.code() == error_code_not_committed)
					++self->conflicts;
				
				wait(tr.onError(e));
				++self->retries;
			}
			// reset all the operations' counters to 0
			std::fill(perOpCount.begin(), perOpCount.end(), 0);
			tr.reset();
		}
	}

	ACTOR Future<Void> cleanup(Database cx, MakoWorkload* self){
		// clear all data starts with 'mako' in the database
		state std::string keyPrefix(self->KEYPREFIX);
		state ReadYourWritesTransaction tr(cx);

		loop{
			try {
				tr.clear(prefixRange(keyPrefix));
				wait(tr.commit());
				break;
			} catch (Error &e){
				TraceEvent("FailedToCleanData").error(e);
				wait(tr.onError(e));
			}
		}

		return Void();
	}
	ACTOR template<class T>
	static Future<Void> logLatency(Future<T> f,  ContinuousSample<double>* opLatencies){
		state double opBegin = now();
		T value = wait(f);
		opLatencies->addSample(now() - opBegin);
		return Void();
	}

	int64_t getRandomKey(uint64_t rowCount) {
		// TODO: support other distribution like zipf
		return deterministicRandom()->randomInt64(0, rowCount); 
	}
	void parseOperationsSpec() {
		const char *ptr = operationsSpec.c_str();
		int op = 0;
		int rangeop = 0;
		int num;
		int error = 0;

		for (op = 0; op < MAX_OP; op++) {
			operations[op][OP_COUNT] = 0;
			operations[op][OP_RANGE] = 0;
		}

		op = 0;
		while (*ptr) {
			if (strncmp(ptr, "grv", 3) == 0) {
				op = OP_GETREADVERSION;
				ptr += 3;
			} else if (strncmp(ptr, "gr", 2) == 0) {
				op = OP_GETRANGE;
				rangeop = 1;
				ptr += 2;
			} else if (strncmp(ptr, "g", 1) == 0) {
				op = OP_GET;
				ptr++;
			} else if (strncmp(ptr, "sgr", 3) == 0) {
				op = OP_SGETRANGE;
				rangeop = 1;
				ptr += 3;
			} else if (strncmp(ptr, "sg", 2) == 0) {
				op = OP_SGET;
				ptr += 2;
			} else if (strncmp(ptr, "u", 1) == 0) {
				op = OP_UPDATE;
				ptr++;
			} else if (strncmp(ptr, "ir", 2) == 0) {
				op = OP_INSERTRANGE;
				rangeop = 1;
				ptr += 2;
			} else if (strncmp(ptr, "i", 1) == 0) {
				op = OP_INSERT;
				ptr++;
			} else if (strncmp(ptr, "cr", 2) == 0) {
				op = OP_CLEARRANGE;
				rangeop = 1;
				ptr += 2;
			} else if (strncmp(ptr, "c", 1) == 0) {
				op = OP_CLEAR;
				ptr++;
			} else if (strncmp(ptr, "scr", 3) == 0) {
				op = OP_SETCLEARRANGE;
				rangeop = 1;
				ptr += 3;
			} else if (strncmp(ptr, "sc", 2) == 0) {
				op = OP_SETCLEAR;
				ptr += 2;
			} else {
				error = 1;
				break;
			}

			/* count */
			num = 0;
			if ((*ptr < '0') || (*ptr > '9')) {
				num = 1; /* if omitted, set it to 1 */
			} else {
				while ((*ptr >= '0') && (*ptr <= '9')) {
					num = num * 10 + *ptr - '0';
					ptr++;
				}
			}
			/* set count */
			operations[op][OP_COUNT] = num;

			if (rangeop) {
				if (*ptr != ':') {
					error = 1;
					break;
				} else {
					ptr++; /* skip ':' */
					num = 0;
					if ((*ptr < '0') || (*ptr > '9')) {
						error = 1;
						break;
					}
					while ((*ptr >= '0') && (*ptr <= '9')) {
						num = num * 10 + *ptr - '0';
						ptr++;
					}
					/* set range */
					if (num > RANGELIMIT)
						TraceEvent(SevError, "RangeExceedLimit").detail("RangeLimit", RANGELIMIT).detail("Range", num);
					operations[op][OP_RANGE] = num;
				}
			}
			rangeop = 0;
		}

		if (error) {
			TraceEvent(SevError, "InvalidTransactionSpecification").detail("operations", operationsSpec);
		}
	}
};

WorkloadFactory<MakoWorkload> MakoloadFactory("Mako");
