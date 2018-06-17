package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	log "github.com/dihedron/go-log"
)

func init() {
	log.SetLevel(log.DBG)
	log.SetStream(os.Stderr, true)
	log.SetTimeFormat("15:04:05.000")
	log.SetPrintCallerInfo(true)
	log.SetPrintSourceInfo(log.SourceInfoShort)
}

func main() {

	once := flag.Bool("once", false, "whether the instruction should be applied only to the first occurrence")
	flag.Parse()
	processStream(flag.Args(), *once)

	//cmd.Execute()

}

type operation int8

const (
	// OperationReplace replaces the matching line with the user provided text.
	OperationReplace operation = iota
	// OperationPrepend inserts the user provided text before the matching line.
	OperationPrepend
	// OperationAppend inserts the user provided text after the matching line.
	OperationAppend
	// OperationDelete removes the matching line.
	OperationDelete
	// OperationInvalid means that the operation could not be recognised.
	OperationInvalid
)

// processStream is the actual workhorse: it identifies input and output, then
// reads in the input stream one line at a time and applies its pattern matching
// line by line; matching lines are processed and written to the output stream.
func processStream(args []string, once bool) {
	log.Debugf("Apply only once: %t", once)
	for i, arg := range args {
		log.Debugf("args[%d] => %q\n", i, arg)
	}

	input, err := getInput(args)
	if err != nil {
		log.Fatalf("Unable to open input file: %v", err)
	}
	defer input.Close()

	output, err := getOutput(args)
	if err != nil {
		log.Fatalf("Unable to open output file: %v", err)
	}
	defer output.Close()

	op := getOperation(args)
	log.Debugf("Operation: %d", op)

	log.Debugf("Matching against %q", args[2])
	re := regexp.MustCompile(args[2])

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			log.Debugf("Input text %q matches pattern", scanner.Text())
			line := processLine(scanner.Text(), args[0], re)
			switch op {
			case OperationReplace:
				fmt.Fprintf(output, "%s\n", line)
			case OperationPrepend:
				fmt.Fprintf(output, "%s\n", line)
				fmt.Fprintf(output, "%s\n", scanner.Text())
			case OperationAppend:
				fmt.Fprintf(output, "%s\n", scanner.Text())
				fmt.Fprintf(output, "%s\n", line)
			case OperationDelete:
			}

			// if placeholders.MatchString(args[2]) {
			// 	log.Debugf("Replacement text has bindings")
			// 	// TODO: find all capturing groups in scanner.Text(), then use them to
			// 	// bind the replacement arguments; this processing is common to all
			// 	// matching methods so it should be moved to its own method.
			// 	for _, indexes := range placeholders.FindAllStringSubmatchIndex(args[2], -1) {
			// 		index, _ := strconv.Atoi(args[2][indexes[2]:indexes[3]])
			// 		log.Debugf("Match: %q (%d) from %d to %d", args[2][indexes[0]:indexes[1]], index, indexes[0], indexes[1])
			// 	}
			// 	//matches := re.FindStringSubmatch(scanner.Text())
			// } else {
			// 	log.Debugf("Replacing text %q with %q\n", scanner.Text(), args[2])
			// 	fmt.Fprintf(output, "%s\n", args[2])
			// }
		} else {
			log.Debugf("Keeping text %q\n", scanner.Text())
			fmt.Fprintf(output, "%s\n", scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading text: %v", err)
	}

}

// getInput returns the input Reader to use; if a filename argument is provided,
// open the file to read from it, otherwise return STDIN; the Reader must be
// closed by the method's caller.
func getInput(args []string) (*os.File, error) {
	if len(args) > 3 && args[3] != "" {
		log.Debugf("Reading text from input file: %q", args[3])
		return os.Open(args[3])
	}
	return os.Stdin, nil
}

// getOutput returns the output Writer to use; if a filename argument is provided,
// open the file to write to it, otherwise return STDOUT; the Writer must be
// closed by the method's caller.
func getOutput(args []string) (*os.File, error) {
	if len(args) > 4 && args[4] != "" {
		log.Debugf("Writing text to output file: %q", args[4])
		return os.Create(args[4])
	}
	return os.Stdout, nil
}

// getOperation decodes the requested operation using the clause, according to the
// command usage; fuzzy matching (see github.com/sahilm/fuzzy) may be introduced
// later on once the product is sufficiently stable.
func getOperation(args []string) operation {
	if args[1] == "where" || args[1] == "wherever" {
		if args[0] == "-" {
			return OperationDelete
		}
		return OperationReplace
	}
	if args[1] == "before" {
		return OperationPrepend
	}
	if args[1] == "after" {
		return OperationAppend
	}
	log.Fatalf("Unknown clause: %q; valid values include 'where', 'wherever', after' and 'before'")
	return OperationInvalid
}

var anchors = regexp.MustCompile(`(?:\{(\d+)\})`)

func processLine(original string, replacement string, re *regexp.Regexp) string {
	if anchors.MatchString(replacement) {
		log.Debugf("Replacement text requires binding\n")
		// TODO: find all capturing groups in scanner.Text(), then use them to
		// bind the replacement arguments; this processing is common to all
		// matching methods so it should be moved to its own method.
		matches := re.FindStringSubmatch(original)
		if len(matches) == 0 {
			log.Fatalf("Invalid number of bindings: %d\n", len(matches))
		}

		bindings := []string{}
		for i, match := range matches {
			log.Debugf("Match[%d] => %q\n", i, match)
			bindings = append(bindings, match)
		}

		buffer := ""
		cursor := 0
		// lengthening := 0
		for _, indexes := range anchors.FindAllStringSubmatchIndex(replacement, -1) {
			index, _ := strconv.Atoi(replacement[indexes[2]:indexes[3]])
			if index > len(bindings) {
				var buffer bytes.Buffer
				buffer.WriteString(fmt.Sprintf("Invalid binding index %d, current binding values are:\n", index))
				for i, binding := range bindings {
					buffer.WriteString(fmt.Sprintf(" {%d} => %q\n", i, binding))
				}
				log.Fatalln(buffer.String())
			}
			log.Debugf("Match: %q (%d) from %d to %d", replacement[indexes[0]:indexes[1]], index, indexes[0], indexes[1])
			buffer = buffer + replacement[cursor:indexes[0]] + bindings[index]
			cursor = indexes[1]
			log.Debugf("Current temporary buffer: %q", buffer)
		}
		buffer = buffer + replacement[cursor:]
		log.Debugf("Current temporary buffer: %q", buffer)
		return buffer
	} else {
		log.Debugf("Replacing text %q with %q\n", original, replacement)
		return replacement
	}
}
