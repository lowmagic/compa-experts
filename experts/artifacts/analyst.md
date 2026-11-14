# Unfamiliar artifacts

For a directory of files whose layout nobody has described: an archive unpacked by hand, a product's own diagnostic
bundle, a folder of logs. It knows no product. Where an expert for the product is installed, load it; this one is for
finding out what is there.

## Method

Work from the outside in, and let each step choose the next.

1. **Lay it out.** `generic.files.tree` gives the directories to three levels, and `generic.files.ls` what one
   directory holds with each entry's size and modification time. Sizes and times are evidence: an empty file, a file
   far larger than its neighbours, or one last written at the time of the incident says where to look.
2. **Find by name.** `generic.files.find` lists the files whose names match a glob, such as `*.log` or `*.conf`. A name
   that is absent is a fact about the listing, not about the system the files came from: say what was listed and to
   what depth.
3. **Count before reading.** `generic.files.count` groups a log's lines and gives each group's count with its first and
   last occurrence. It tells a message printed once from one printed ten thousand times, and when each began and
   stopped, for one call.
4. **Search.** `generic.files.grep` returns the lines matching a pattern with their file and line number. Start with
   the exact string from the report: an error message, an identifier, a time. Add context only around the matches
   that matter.
5. **Sample.** `generic.files.head` and `generic.files.tail` read a file's first and last lines. A file's head often
   says what wrote it and in which format; a log's tail says how it ended.

## Reading what comes back

- A listing shows what was collected, never what was running. A configuration file shows what was written, not what
  was applied.
- Two files are about the same thing only when something in them says so: a host name, an identifier, a time in the
  same zone.
- Times without a zone are not comparable with times that have one until the zone is found.
- A search that returns nothing is evidence only for the files and the pattern it was given. Name both.
- A result marked truncated or stopped is partial. Narrow the path or the pattern and run it again before reasoning
  from it.

## When to stop

When the files identify a product or a subsystem, say which, with the line that shows it, and load or consult its
expert from the index. When none is installed, say so, and go on with these tools and the engine's own.
