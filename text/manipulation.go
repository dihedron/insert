package text

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	log "github.com/dihedron/go-log"
	"github.com/spf13/cobra"
)

func Copy(cmd *cobra.Command, args []string) {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintf(os.Stdout, "%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading text: %v", err)
	}
}

func Replace(cmd *cobra.Command, args []string) {
	if args[1] != "with" {
		log.Fatalf("Error: 'with' clause was not specified")
	}

	// if a filename argument is provided, read from file, otherwise it's STDIN
	var input *os.File
	if len(args) > 3 && args[3] != "" {
		log.Debugf("Reading text from input file: %q", args[3])
		var err error
		input, err = os.Open(args[3])
		if err != nil {
			log.Fatalf("Unable to open input file: %v", err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	// if a filename argument is provided, write to file, otherwise it's STDOUT
	var output *os.File
	if len(args) > 4 && args[4] != "" {
		log.Debugf("Writing text to output file: %q", args[4])
		var err error
		output, err = os.Create(args[4])
		if err != nil {
			log.Fatalf("Unable to open output file: %v", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	log.Debugf("Matching against %q", args[0])
	re := regexp.MustCompile(args[0])

	re2 := regexp.MustCompile(`(?:\{(\d+)\})`)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			log.Debugf("Input text %q matches pattern", scanner.Text())
			if re2.MatchString(args[2]) {
				log.Debugf("Replacement text has bindings")
				variables := re2.FindAllStringSubmatchIndex(args[2], -1)
				for _, variables2 := range variables {
					index, _ := strconv.Atoi(args[2][variables2[2]:variables2[3]])
					log.Debugf("Match: %q (%d) from %d to %d", args[2][variables2[0]:variables2[1]], index, variables2[0], variables2[1])
					for i, variable := range variables2 {
						log.Debugf(" [%d] => %v", i, variable)
					}
				}
				//matches := re.FindStringSubmatch(scanner.Text())

			} else {
				log.Debugf("Replacing text %q with %q\n", scanner.Text(), args[2])
				fmt.Fprintf(output, "%s\n", args[2])
			}
		} else {
			log.Debugf("Keeping text %q\n", scanner.Text())
			fmt.Fprintf(output, "%s\n", scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading text: %v", err)
	}
}
