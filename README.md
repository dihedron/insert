# zed

```zed``` is a stream editor with a limited subset of ```sed```'s functionalities exposed through a modern CLI; it is designed to be used as a shell filter, with readability in mind.

## Usage

```zed``` is meant to be used like this:

```bash
$ > cat infile.txt | 
        zed replace "^a\s*line$" with "another line" --everywhere | 
        zed after "^yet\s+another\s*line$" insert "some text" | 
        zed before "^and this is the (last)$" insert "${1} but not least" --once 
        > out.txt 
```
