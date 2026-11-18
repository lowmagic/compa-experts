# Failed connections in a service's logs

For a service that stops working when it cannot reach another: a worker losing its database, a client its server. It
reads what the service's own logs say about the connections it tried. It knows no product; where an expert for the
service or its peer is installed, load it.

## Method

1. **Find the failures.** `generic.connections.failures` returns the lines of one log that report a failed connection,
   with their line numbers. It reports unavailable when the file is not among the artifacts: that is the file missing
   from the case, not a log without failures.
2. **Count them.** `generic.files.count` with the same file and a pattern such as `refused` gives how often each
   failure repeats and its first and last occurrence. Set those times beside when the symptom began and when it
   stopped.
3. **Read how it ended.** `generic.files.tail` gives the log's last lines: whether the service gave up, retried or
   recovered after the failures.
4. **Search for the peer.** `generic.files.grep` for the address or name the failures give finds what else the
   artifacts say about that peer.

## What the failures say

The wording below is the C library's for the error a connection attempt returns (Linux man-pages `connect(2)` and
`errno(3)`); a service may print it, or its own words for it.

- **Connection refused** (`ECONNREFUSED`): the attempt found no one listening on the remote address.
- **Connection timed out** (`ETIMEDOUT`): the attempt timed out; the server may have been too busy to accept new
  connections.
- **Connection reset** (`ECONNRESET`): a connection was reset.
- **Network is unreachable** (`ENETUNREACH`) and **no route to host** (`EHOSTUNREACH`): the address could not be
  reached at all.

A client's log says what the client was told, not why: that nothing listened at the address does not say whether the
peer was down, restarting or listening elsewhere. Say what the log shows, with its lines, and what would tell the
explanations apart: the peer's own logs for the same times, or the customer checking that the peer listens at the
address the client uses.

## When to stop

When the failures name a peer whose expert is installed, load or consult it with the lines that show the failures. When
none is installed, say what the client's log establishes and what only the peer's side can.
