# zed

```zed``` is a stream editor with a limited subset of ```sed```'s functionalities exposed through a modern CLI; it is designed to be used as a shell filter, with readability in mind.

## Usage

```zed``` is meant to be used like this:

```bash
$ > cat infile.txt | 
        put "another line" where "^a\s*line$" --everywhere | 
        put "some text" after "^yet\s+another\s*line$" | 
        put "${1} but not least" before "^and this is the (last)$" --once 
        put none where "^a\s+line\s*to drop$"  
        > out.txt 
```
