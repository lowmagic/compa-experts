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
`errno(3)`); a service may print it, or its own words for it. Each is what the client was told.

- **Connection refused** (`ECONNREFUSED`): in `connect(2)`'s words, the attempt found no one listening on the remote
  address.
- **Connection timed out** (`ETIMEDOUT`): the attempt timed out before a connection was made.
- **Connection reset** (`ECONNRESET`): a connection was reset.
- **Network is unreachable** (`ENETUNREACH`) and **no route to host** (`EHOSTUNREACH`): the client was told the
  network, or the host, the address is on was unreachable.

A client's log says what the client was told, never why: name no cause for it. Say what the log shows, with its lines,
the address each names and the times they span. What bears on the peer's side at those times is the peer's own logs
for the same times, and whether the peer listened then at the address and port the client used: a question to ask,
which the client's log does not answer.

## When to stop

When the failures name a peer whose expert is installed, load or consult it with the lines that show the failures. When
none is installed, say what the client's log establishes and what only the peer's side can.
