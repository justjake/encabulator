package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/justjake/encabulator/task"
	"github.com/justjake/encabulator/unison"
	"log"
	"os/exec"
	"path"
	"regexp"
	"time"
)

var unisonDelim = regexp.MustCompile("\r\n|\n|\r")

func init() {
	unisonDelim.Longest()
}

func main() {

	manager := unison.Manager()
	manifest := unison.OsManifest()
	for _, req := range manifest {
		err := manager.Stage(req, path.Base(req))
		if err != nil {
			panic(err)
		}
	}

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatalf("Requires arguments SOURCE and DESTINATION")
	}

	fmt.Printf("SOURCE: %q\nDESTINATION: %q\n", args[0], args[1])
	// TODO: flag for repeat
	cmd, err := unison.Unison(args[0], args[1], true, args[2:]...)

	if err != nil {
		log.Fatal(err)
	}

	runForever(cmd.Build())
}

func runForever(cmd *exec.Cmd) {
	splitter := splitByRegexp(unisonDelim)
	supervisor := task.MakeSupervisor(1, time.Second)
	t, err := task.Spawn(cmd, splitter)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		event := <-t.Output
		t, err = supervisor.HandleEvent(event)
		if err != nil {
			log.Fatalln(err)
		}
		if output, ok := event.Payload.(*task.Output); ok {
			log.Printf("%q", output.Chunk)
		}
	}
}

func splitByRegexp(delim *regexp.Regexp) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if loc := delim.FindIndex(data); loc != nil {
			end := loc[1]
			// If the match occurs the end of data, there's a possibility that a
			// longer match is possible if we had more bytes. In that case, just
			// request more data.
			//
			// note that a regexp like /$/ that always matches the end, would cause a
			// panic.
			//
			// I disabled this idea after realizing that in a real-time streamming
			// system, a delimiter regexp almost always occurs at the end of the
			// match, so this deplays the processing of data whole tick! Bad idea.
			/*
				if !atEOF && end == len(data)-1 {
					log.Println("match at EOS")
					return 0, nil, nil
				}
			*/

			return end, data[0:end], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
}

func splitBy(delims []byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		minDistance := -1
		// TODO: Investigate this rather naïve strategy. Benchmark this against
		// golang's regexp stuff, which wouldn't look at the same byte twice.
		for _, delim := range delims {
			if i := bytes.IndexByte(data, delim); i >= 0 {
				if i <= minDistance || minDistance == -1 {
					// we've seen a delim
					minDistance = i
				}
				// We have a full newline-terminated line.
				//return i + 1, dropCR(data[0:i]), nil
			}
		}

		if minDistance != -1 {
			return minDistance + 1, data[0:minDistance], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
}
