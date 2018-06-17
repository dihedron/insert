# put - Yet another stream editor

```put``` is a stream editor with a limited subset of ```sed```'s functionalities exposed through a modern CLI; it is designed to be used as a shell filter, with readability in mind.

## Usage

```put``` is meant to be used like this:

```bash
$ > cat infile.txt | 
        put "some replacement text" where "^a\s*line$" | 
        put --once "{1} but not least" before "^and this is the (last)$" |
        put "some text" after "^yet\s+another\s*line$" |         
        put - where "^a\s+line\s*to drop$"  
        > out.txt 
```

It provides 4 types of operations:
1. replace lines matching a pattern with other text (```where``` clause with replacement text);
2. add a line before each line matching a pattern (```before``` clause);
3. add a line after each line matching a pattern (```after``` clause);
4. drop lines matching a pattern (```where``` clause with ```-``` as replacement text).

The ```--once``` flags indicates that the operaton should occur only on the first occurrence of the match; if omitted, each matching line is "edited" as instructed.

Replacement text can include substitution anchors, such as ```{1}``` above; it will be substituted with the value of the first capturing group in the pattern (the zero-th being the whole expression unless specified otherwise). 

## Example

As an example let's take an ```/etc/hosts```; say you want to add the host name on the ```localhost``` line (the one starting with ```127.0.0.1```) to prevent complaining by your ```sudo``` commands. The following sequence copies the current line to a comment (first invocation of ```put```), then replaces whatever is after the ```localhost``` word with the current hostname:

```bash
$ > cat hosts | 
        put "# {0}" before "^(127\.0\.0\.1\s+localhost).*" | 
        put "{1} $(hostname)" where "^(127\.0\.0\.1\s+localhost).*" 
        > hosts2
```

## Debugging

If you want to see what the command is doing internally, simply run it with the ```PUT_DEBUG``` environment variable set to one of ```debug```, ```info```, ```warning``` or ```error```, e.g. as follows:

```bash
$ > PUT_DEBUG=debug put [args] < /etc/hosts
```

This can be hepful if you need to see what are the available bindings for your regular expression (```{0}```, ```{1}```...) so you can debug it. 

## Suggestions and contributions

... are very welcome: please open an Issue to give your feedback or to request a pull.
