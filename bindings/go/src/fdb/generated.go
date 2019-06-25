/*
 * generated.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// DO NOT EDIT THIS FILE BY HAND. This file was generated using
// translate_fdb_options.go, part of the FoundationDB repository, and a copy of
// the fdb.options file (installed as part of the FoundationDB client, typically
// found as /usr/include/foundationdb/fdb.options).

// To regenerate this file, from the top level of a FoundationDB repository
// checkout, run:
// $ go run bindings/go/src/_util/translate_fdb_options.go < fdbclient/vexillographer/fdb.options > bindings/go/src/fdb/generated.go

package fdb

import (
	"encoding/binary"
)

func int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(i))
	return buf
}

// Deprecated
//
// Parameter: IP:PORT
func (o NetworkOptions) SetLocalAddress(param string) error {
	return o.setOpt(10, []byte(param))
}

// enable the object serializer for network communication
//
// Parameter: 0 is false, every other value is true
func (o NetworkOptions) SetUseObjectSerializer(param int64) error {
	return o.setOpt(11, int64ToBytes(param))
}

// Deprecated
//
// Parameter: path to cluster file
func (o NetworkOptions) SetClusterFile(param string) error {
	return o.setOpt(20, []byte(param))
}

// Enables trace output to a file in a directory of the clients choosing
//
// Parameter: path to output directory (or NULL for current working directory)
func (o NetworkOptions) SetTraceEnable(param string) error {
	return o.setOpt(30, []byte(param))
}

// Sets the maximum size in bytes of a single trace output file. This value should be in the range ``[0, INT64_MAX]``. If the value is set to 0, there is no limit on individual file size. The default is a maximum size of 10,485,760 bytes.
//
// Parameter: max size of a single trace output file
func (o NetworkOptions) SetTraceRollSize(param int64) error {
	return o.setOpt(31, int64ToBytes(param))
}

// Sets the maximum size of all the trace output files put together. This value should be in the range ``[0, INT64_MAX]``. If the value is set to 0, there is no limit on the total size of the files. The default is a maximum size of 104,857,600 bytes. If the default roll size is used, this means that a maximum of 10 trace files will be written at a time.
//
// Parameter: max total size of trace files
func (o NetworkOptions) SetTraceMaxLogsSize(param int64) error {
	return o.setOpt(32, int64ToBytes(param))
}

// Sets the 'LogGroup' attribute with the specified value for all events in the trace output files. The default log group is 'default'.
//
// Parameter: value of the LogGroup attribute
func (o NetworkOptions) SetTraceLogGroup(param string) error {
	return o.setOpt(33, []byte(param))
}

// Select the format of the log files. xml (the default) and json are supported.
//
// Parameter: Format of trace files
func (o NetworkOptions) SetTraceFormat(param string) error {
	return o.setOpt(34, []byte(param))
}

// Set internal tuning or debugging knobs
//
// Parameter: knob_name=knob_value
func (o NetworkOptions) SetKnob(param string) error {
	return o.setOpt(40, []byte(param))
}

// Deprecated
//
// Parameter: file path or linker-resolved name
func (o NetworkOptions) SetTLSPlugin(param string) error {
	return o.setOpt(41, []byte(param))
}

// Set the certificate chain
//
// Parameter: certificates
func (o NetworkOptions) SetTLSCertBytes(param []byte) error {
	return o.setOpt(42, param)
}

// Set the file from which to load the certificate chain
//
// Parameter: file path
func (o NetworkOptions) SetTLSCertPath(param string) error {
	return o.setOpt(43, []byte(param))
}

// Set the private key corresponding to your own certificate
//
// Parameter: key
func (o NetworkOptions) SetTLSKeyBytes(param []byte) error {
	return o.setOpt(45, param)
}

// Set the file from which to load the private key corresponding to your own certificate
//
// Parameter: file path
func (o NetworkOptions) SetTLSKeyPath(param string) error {
	return o.setOpt(46, []byte(param))
}

// Set the peer certificate field verification criteria
//
// Parameter: verification pattern
func (o NetworkOptions) SetTLSVerifyPeers(param []byte) error {
	return o.setOpt(47, param)
}

// Not yet implemented.
func (o NetworkOptions) SetBuggifyEnable() error {
	return o.setOpt(48, nil)
}

// Not yet implemented.
func (o NetworkOptions) SetBuggifyDisable() error {
	return o.setOpt(49, nil)
}

// Set the probability of a BUGGIFY section being active for the current execution.  Only applies to code paths first traversed AFTER this option is changed.
//
// Parameter: probability expressed as a percentage between 0 and 100
func (o NetworkOptions) SetBuggifySectionActivatedProbability(param int64) error {
	return o.setOpt(50, int64ToBytes(param))
}

// Set the probability of an active BUGGIFY section being fired
//
// Parameter: probability expressed as a percentage between 0 and 100
func (o NetworkOptions) SetBuggifySectionFiredProbability(param int64) error {
	return o.setOpt(51, int64ToBytes(param))
}

// Set the ca bundle
//
// Parameter: ca bundle
func (o NetworkOptions) SetTLSCaBytes(param []byte) error {
	return o.setOpt(52, param)
}

// Set the file from which to load the certificate authority bundle
//
// Parameter: file path
func (o NetworkOptions) SetTLSCaPath(param string) error {
	return o.setOpt(53, []byte(param))
}

// Set the passphrase for encrypted private key. Password should be set before setting the key for the password to be used.
//
// Parameter: key passphrase
func (o NetworkOptions) SetTLSPassword(param string) error {
	return o.setOpt(54, []byte(param))
}

// Disables the multi-version client API and instead uses the local client directly. Must be set before setting up the network.
func (o NetworkOptions) SetDisableMultiVersionClientApi() error {
	return o.setOpt(60, nil)
}

// If set, callbacks from external client libraries can be called from threads created by the FoundationDB client library. Otherwise, callbacks will be called from either the thread used to add the callback or the network thread. Setting this option can improve performance when connected using an external client, but may not be safe to use in all environments. Must be set before setting up the network. WARNING: This feature is considered experimental at this time.
func (o NetworkOptions) SetCallbacksOnExternalThreads() error {
	return o.setOpt(61, nil)
}

// Adds an external client library for use by the multi-version client API. Must be set before setting up the network.
//
// Parameter: path to client library
func (o NetworkOptions) SetExternalClientLibrary(param string) error {
	return o.setOpt(62, []byte(param))
}

// Searches the specified path for dynamic libraries and adds them to the list of client libraries for use by the multi-version client API. Must be set before setting up the network.
//
// Parameter: path to directory containing client libraries
func (o NetworkOptions) SetExternalClientDirectory(param string) error {
	return o.setOpt(63, []byte(param))
}

// Prevents connections through the local client, allowing only connections through externally loaded client libraries. Intended primarily for testing.
func (o NetworkOptions) SetDisableLocalClient() error {
	return o.setOpt(64, nil)
}

// Disables logging of client statistics, such as sampled transaction activity.
func (o NetworkOptions) SetDisableClientStatisticsLogging() error {
	return o.setOpt(70, nil)
}

// Enables debugging feature to perform slow task profiling. Requires trace logging to be enabled. WARNING: this feature is not recommended for use in production.
func (o NetworkOptions) SetEnableSlowTaskProfiling() error {
	return o.setOpt(71, nil)
}

// Set the size of the client location cache. Raising this value can boost performance in very large databases where clients access data in a near-random pattern. Defaults to 100000.
//
// Parameter: Max location cache entries
func (o DatabaseOptions) SetLocationCacheSize(param int64) error {
	return o.setOpt(10, int64ToBytes(param))
}

// Set the maximum number of watches allowed to be outstanding on a database connection. Increasing this number could result in increased resource usage. Reducing this number will not cancel any outstanding watches. Defaults to 10000 and cannot be larger than 1000000.
//
// Parameter: Max outstanding watches
func (o DatabaseOptions) SetMaxWatches(param int64) error {
	return o.setOpt(20, int64ToBytes(param))
}

// Specify the machine ID that was passed to fdbserver processes running on the same machine as this client, for better location-aware load balancing.
//
// Parameter: Hexadecimal ID
func (o DatabaseOptions) SetMachineId(param string) error {
	return o.setOpt(21, []byte(param))
}

// Specify the datacenter ID that was passed to fdbserver processes running in the same datacenter as this client, for better location-aware load balancing.
//
// Parameter: Hexadecimal ID
func (o DatabaseOptions) SetDatacenterId(param string) error {
	return o.setOpt(22, []byte(param))
}

// Set a timeout in milliseconds which, when elapsed, will cause each transaction automatically to be cancelled. This sets the ``timeout`` option of each transaction created by this database. See the transaction option description for more information. Using this option requires that the API version is 610 or higher.
//
// Parameter: value in milliseconds of timeout
func (o DatabaseOptions) SetTransactionTimeout(param int64) error {
	return o.setOpt(500, int64ToBytes(param))
}

// Set a timeout in milliseconds which, when elapsed, will cause a transaction automatically to be cancelled. This sets the ``retry_limit`` option of each transaction created by this database. See the transaction option description for more information.
//
// Parameter: number of times to retry
func (o DatabaseOptions) SetTransactionRetryLimit(param int64) error {
	return o.setOpt(501, int64ToBytes(param))
}

// Set the maximum amount of backoff delay incurred in the call to ``onError`` if the error is retryable. This sets the ``max_retry_delay`` option of each transaction created by this database. See the transaction option description for more information.
//
// Parameter: value in milliseconds of maximum delay
func (o DatabaseOptions) SetTransactionMaxRetryDelay(param int64) error {
	return o.setOpt(502, int64ToBytes(param))
}

// Set the maximum transaction size which, if exceeded, will cause the transaction to be cancelled. Default to 10,000,000 bytes.
//
// Parameter: value in bytes
func (o DatabaseOptions) SetTransactionSizeLimit(param int64) error {
	return o.setOpt(503, int64ToBytes(param))
}

// Snapshot read operations will see the results of writes done in the same transaction. This is the default behavior.
func (o DatabaseOptions) SetSnapshotRywEnable() error {
	return o.setOpt(26, nil)
}

// Snapshot read operations will not see the results of writes done in the same transaction. This was the default behavior prior to API version 300.
func (o DatabaseOptions) SetSnapshotRywDisable() error {
	return o.setOpt(27, nil)
}

// The transaction, if not self-conflicting, may be committed a second time after commit succeeds, in the event of a fault
func (o TransactionOptions) SetCausalWriteRisky() error {
	return o.setOpt(10, nil)
}

// The read version will be committed, and usually will be the latest committed, but might not be the latest committed in the event of a fault or partition
func (o TransactionOptions) SetCausalReadRisky() error {
	return o.setOpt(20, nil)
}

// Not yet implemented.
func (o TransactionOptions) SetCausalReadDisable() error {
	return o.setOpt(21, nil)
}

// The next write performed on this transaction will not generate a write conflict range. As a result, other transactions which read the key(s) being modified by the next write will not conflict with this transaction. Care needs to be taken when using this option on a transaction that is shared between multiple threads. When setting this option, write conflict ranges will be disabled on the next write operation, regardless of what thread it is on.
func (o TransactionOptions) SetNextWriteNoWriteConflictRange() error {
	return o.setOpt(30, nil)
}

// Reads performed by a transaction will not see any prior mutations that occured in that transaction, instead seeing the value which was in the database at the transaction's read version. This option may provide a small performance benefit for the client, but also disables a number of client-side optimizations which are beneficial for transactions which tend to read and write the same keys within a single transaction.
func (o TransactionOptions) SetReadYourWritesDisable() error {
	return o.setOpt(51, nil)
}

// Deprecated
func (o TransactionOptions) SetReadAheadDisable() error {
	return o.setOpt(52, nil)
}

// Not yet implemented.
func (o TransactionOptions) SetDurabilityDatacenter() error {
	return o.setOpt(110, nil)
}

// Not yet implemented.
func (o TransactionOptions) SetDurabilityRisky() error {
	return o.setOpt(120, nil)
}

// Deprecated
func (o TransactionOptions) SetDurabilityDevNullIsWebScale() error {
	return o.setOpt(130, nil)
}

// Specifies that this transaction should be treated as highest priority and that lower priority transactions should block behind this one. Use is discouraged outside of low-level tools
func (o TransactionOptions) SetPrioritySystemImmediate() error {
	return o.setOpt(200, nil)
}

// Specifies that this transaction should be treated as low priority and that default priority transactions will be processed first. Batch priority transactions will also be throttled at load levels smaller than for other types of transactions and may be fully cut off in the event of machine failures. Useful for doing batch work simultaneously with latency-sensitive work
func (o TransactionOptions) SetPriorityBatch() error {
	return o.setOpt(201, nil)
}

// This is a write-only transaction which sets the initial configuration. This option is designed for use by database system tools only.
func (o TransactionOptions) SetInitializeNewDatabase() error {
	return o.setOpt(300, nil)
}

// Allows this transaction to read and modify system keys (those that start with the byte 0xFF)
func (o TransactionOptions) SetAccessSystemKeys() error {
	return o.setOpt(301, nil)
}

// Allows this transaction to read system keys (those that start with the byte 0xFF)
func (o TransactionOptions) SetReadSystemKeys() error {
	return o.setOpt(302, nil)
}

// Not yet implemented.
func (o TransactionOptions) SetDebugRetryLogging(param string) error {
	return o.setOpt(401, []byte(param))
}

// Deprecated
//
// Parameter: String identifier to be used in the logs when tracing this transaction. The identifier must not exceed 100 characters.
func (o TransactionOptions) SetTransactionLoggingEnable(param string) error {
	return o.setOpt(402, []byte(param))
}

// Sets a client provided identifier for the transaction that will be used in scenarios like tracing or profiling. Client trace logging or transaction profiling must be separately enabled.
//
// Parameter: String identifier to be used when tracing or profiling this transaction. The identifier must not exceed 100 characters.
func (o TransactionOptions) SetDebugTransactionIdentifier(param string) error {
	return o.setOpt(403, []byte(param))
}

// Enables tracing for this transaction and logs results to the client trace logs. The DEBUG_TRANSACTION_IDENTIFIER option must be set before using this option, and client trace logging must be enabled and to get log output.
func (o TransactionOptions) SetLogTransaction() error {
	return o.setOpt(404, nil)
}

// Set a timeout in milliseconds which, when elapsed, will cause the transaction automatically to be cancelled. Valid parameter values are ``[0, INT_MAX]``. If set to 0, will disable all timeouts. All pending and any future uses of the transaction will throw an exception. The transaction can be used again after it is reset. Prior to API version 610, like all other transaction options, the timeout must be reset after a call to ``onError``. If the API version is 610 or greater, the timeout is not reset after an ``onError`` call. This allows the user to specify a longer timeout on specific transactions than the default timeout specified through the ``transaction_timeout`` database option without the shorter database timeout cancelling transactions that encounter a retryable error. Note that at all API versions, it is safe and legal to set the timeout each time the transaction begins, so most code written assuming the older behavior can be upgraded to the newer behavior without requiring any modification, and the caller is not required to implement special logic in retry loops to only conditionally set this option.
//
// Parameter: value in milliseconds of timeout
func (o TransactionOptions) SetTimeout(param int64) error {
	return o.setOpt(500, int64ToBytes(param))
}

// Set a maximum number of retries after which additional calls to ``onError`` will throw the most recently seen error code. Valid parameter values are ``[-1, INT_MAX]``. If set to -1, will disable the retry limit. Prior to API version 610, like all other transaction options, the retry limit must be reset after a call to ``onError``. If the API version is 610 or greater, the retry limit is not reset after an ``onError`` call. Note that at all API versions, it is safe and legal to set the retry limit each time the transaction begins, so most code written assuming the older behavior can be upgraded to the newer behavior without requiring any modification, and the caller is not required to implement special logic in retry loops to only conditionally set this option.
//
// Parameter: number of times to retry
func (o TransactionOptions) SetRetryLimit(param int64) error {
	return o.setOpt(501, int64ToBytes(param))
}

// Set the maximum amount of backoff delay incurred in the call to ``onError`` if the error is retryable. Defaults to 1000 ms. Valid parameter values are ``[0, INT_MAX]``. If the maximum retry delay is less than the current retry delay of the transaction, then the current retry delay will be clamped to the maximum retry delay. Prior to API version 610, like all other transaction options, the maximum retry delay must be reset after a call to ``onError``. If the API version is 610 or greater, the retry limit is not reset after an ``onError`` call. Note that at all API versions, it is safe and legal to set the maximum retry delay each time the transaction begins, so most code written assuming the older behavior can be upgraded to the newer behavior without requiring any modification, and the caller is not required to implement special logic in retry loops to only conditionally set this option.
//
// Parameter: value in milliseconds of maximum delay
func (o TransactionOptions) SetMaxRetryDelay(param int64) error {
	return o.setOpt(502, int64ToBytes(param))
}

// Set the maximum transaction size which, if exceeded, will cause the transaction to be cancelled. Valid parameter values are ``[32, 10,000,000]```.
//
// Parameter: value in bytes
func (o TransactionOptions) SetSizeLimit(param int64) error {
	return o.setOpt(503, int64ToBytes(param))
}

// Snapshot read operations will see the results of writes done in the same transaction. This is the default behavior.
func (o TransactionOptions) SetSnapshotRywEnable() error {
	return o.setOpt(600, nil)
}

// Snapshot read operations will not see the results of writes done in the same transaction. This was the default behavior prior to API version 300.
func (o TransactionOptions) SetSnapshotRywDisable() error {
	return o.setOpt(601, nil)
}

// The transaction can read and write to locked databases, and is resposible for checking that it took the lock.
func (o TransactionOptions) SetLockAware() error {
	return o.setOpt(700, nil)
}

// By default, operations that are performed on a transaction while it is being committed will not only fail themselves, but they will attempt to fail other in-flight operations (such as the commit) as well. This behavior is intended to help developers discover situations where operations could be unintentionally executed after the transaction has been reset. Setting this option removes that protection, causing only the offending operation to fail.
func (o TransactionOptions) SetUsedDuringCommitProtectionDisable() error {
	return o.setOpt(701, nil)
}

// The transaction can read from locked databases.
func (o TransactionOptions) SetReadLockAware() error {
	return o.setOpt(702, nil)
}

// This option should only be used by tools which change the database configuration.
func (o TransactionOptions) SetUseProvisionalProxies() error {
	return o.setOpt(711, nil)
}

type StreamingMode int

const (

	// Client intends to consume the entire range and would like it all
	// transferred as early as possible.
	StreamingModeWantAll StreamingMode = -1

	// The default. The client doesn't know how much of the range it is likely
	// to used and wants different performance concerns to be balanced. Only a
	// small portion of data is transferred to the client initially (in order to
	// minimize costs if the client doesn't read the entire range), and as the
	// caller iterates over more items in the range larger batches will be
	// transferred in order to minimize latency.
	StreamingModeIterator StreamingMode = 0

	// Infrequently used. The client has passed a specific row limit and wants
	// that many rows delivered in a single batch. Because of iterator operation
	// in client drivers make request batches transparent to the user, consider
	// ``WANT_ALL`` StreamingMode instead. A row limit must be specified if this
	// mode is used.
	StreamingModeExact StreamingMode = 1

	// Infrequently used. Transfer data in batches small enough to not be much
	// more expensive than reading individual rows, to minimize cost if
	// iteration stops early.
	StreamingModeSmall StreamingMode = 2

	// Infrequently used. Transfer data in batches sized in between small and
	// large.
	StreamingModeMedium StreamingMode = 3

	// Infrequently used. Transfer data in batches large enough to be, in a
	// high-concurrency environment, nearly as efficient as possible. If the
	// client stops iteration early, some disk and network bandwidth may be
	// wasted. The batch size may still be too small to allow a single client to
	// get high throughput from the database, so if that is what you need
	// consider the SERIAL StreamingMode.
	StreamingModeLarge StreamingMode = 4

	// Transfer data in batches large enough that an individual client can get
	// reasonable read bandwidth from the database. If the client stops
	// iteration early, considerable disk and network bandwidth may be wasted.
	StreamingModeSerial StreamingMode = 5
)

// Performs an addition of little-endian integers. If the existing value in the database is not present or shorter than ``param``, it is first extended to the length of ``param`` with zero bytes.  If ``param`` is shorter than the existing value in the database, the existing value is truncated to match the length of ``param``. The integers to be added must be stored in a little-endian representation.  They can be signed in two's complement representation or unsigned. You can add to an integer at a known offset in the value by prepending the appropriate number of zero bytes to ``param`` and padding with zero bytes to match the length of the value. However, this offset technique requires that you know the addition will not cause the integer field within the value to overflow.
func (t Transaction) Add(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 2)
}

// Deprecated
func (t Transaction) And(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 6)
}

// Performs a bitwise ``and`` operation.  If the existing value in the database is not present, then ``param`` is stored in the database. If the existing value in the database is shorter than ``param``, it is first extended to the length of ``param`` with zero bytes.  If ``param`` is shorter than the existing value in the database, the existing value is truncated to match the length of ``param``.
func (t Transaction) BitAnd(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 6)
}

// Deprecated
func (t Transaction) Or(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 7)
}

// Performs a bitwise ``or`` operation.  If the existing value in the database is not present or shorter than ``param``, it is first extended to the length of ``param`` with zero bytes.  If ``param`` is shorter than the existing value in the database, the existing value is truncated to match the length of ``param``.
func (t Transaction) BitOr(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 7)
}

// Deprecated
func (t Transaction) Xor(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 8)
}

// Performs a bitwise ``xor`` operation.  If the existing value in the database is not present or shorter than ``param``, it is first extended to the length of ``param`` with zero bytes.  If ``param`` is shorter than the existing value in the database, the existing value is truncated to match the length of ``param``.
func (t Transaction) BitXor(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 8)
}

// Appends ``param`` to the end of the existing value already in the database at the given key (or creates the key and sets the value to ``param`` if the key is empty). This will only append the value if the final concatenated value size is less than or equal to the maximum value size (i.e., if it fits). WARNING: No error is surfaced back to the user if the final value is too large because the mutation will not be applied until after the transaction has been committed. Therefore, it is only safe to use this mutation type if one can guarantee that one will keep the total value size under the maximum size.
func (t Transaction) AppendIfFits(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 9)
}

// Performs a little-endian comparison of byte strings. If the existing value in the database is not present or shorter than ``param``, it is first extended to the length of ``param`` with zero bytes.  If ``param`` is shorter than the existing value in the database, the existing value is truncated to match the length of ``param``. The larger of the two values is then stored in the database.
func (t Transaction) Max(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 12)
}

// Performs a little-endian comparison of byte strings. If the existing value in the database is not present, then ``param`` is stored in the database. If the existing value in the database is shorter than ``param``, it is first extended to the length of ``param`` with zero bytes.  If ``param`` is shorter than the existing value in the database, the existing value is truncated to match the length of ``param``. The smaller of the two values is then stored in the database.
func (t Transaction) Min(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 13)
}

// Transforms ``key`` using a versionstamp for the transaction. Sets the transformed key in the database to ``param``. The key is transformed by removing the final four bytes from the key and reading those as a little-Endian 32-bit integer to get a position ``pos``. The 10 bytes of the key from ``pos`` to ``pos + 10`` are replaced with the versionstamp of the transaction used. The first byte of the key is position 0. A versionstamp is a 10 byte, unique, monotonically (but not sequentially) increasing value for each committed transaction. The first 8 bytes are the committed version of the database (serialized in big-Endian order). The last 2 bytes are monotonic in the serialization order for transactions. WARNING: At this time, versionstamps are compatible with the Tuple layer only in the Java, Python, and Go bindings. Also, note that prior to API version 520, the offset was computed from only the final two bytes rather than the final four bytes.
func (t Transaction) SetVersionstampedKey(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 14)
}

// Transforms ``param`` using a versionstamp for the transaction. Sets the ``key`` given to the transformed ``param``. The parameter is transformed by removing the final four bytes from ``param`` and reading those as a little-Endian 32-bit integer to get a position ``pos``. The 10 bytes of the parameter from ``pos`` to ``pos + 10`` are replaced with the versionstamp of the transaction used. The first byte of the parameter is position 0. A versionstamp is a 10 byte, unique, monotonically (but not sequentially) increasing value for each committed transaction. The first 8 bytes are the committed version of the database (serialized in big-Endian order). The last 2 bytes are monotonic in the serialization order for transactions. WARNING: At this time, versionstamps are compatible with the Tuple layer only in the Java, Python, and Go bindings. Also, note that prior to API version 520, the versionstamp was always placed at the beginning of the parameter rather than computing an offset.
func (t Transaction) SetVersionstampedValue(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 15)
}

// Performs lexicographic comparison of byte strings. If the existing value in the database is not present, then ``param`` is stored. Otherwise the smaller of the two values is then stored in the database.
func (t Transaction) ByteMin(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 16)
}

// Performs lexicographic comparison of byte strings. If the existing value in the database is not present, then ``param`` is stored. Otherwise the larger of the two values is then stored in the database.
func (t Transaction) ByteMax(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 17)
}

// Performs an atomic ``compare and clear`` operation. If the existing value in the database is equal to the given value, then given key is cleared.
func (t Transaction) CompareAndClear(key KeyConvertible, param []byte) {
	t.atomicOp(key.FDBKey(), param, 20)
}

type conflictRangeType int

const (

	// Used to add a read conflict range
	conflictRangeTypeRead conflictRangeType = 0

	// Used to add a write conflict range
	conflictRangeTypeWrite conflictRangeType = 1
)

type ErrorPredicate int

const (

	// Returns ``true`` if the error indicates the operations in the
	// transactions should be retried because of transient error.
	ErrorPredicateRetryable ErrorPredicate = 50000

	// Returns ``true`` if the error indicates the transaction may have
	// succeeded, though not in a way the system can verify.
	ErrorPredicateMaybeCommitted ErrorPredicate = 50001

	// Returns ``true`` if the error indicates the transaction has not
	// committed, though in a way that can be retried.
	ErrorPredicateRetryableNotCommitted ErrorPredicate = 50002
)
