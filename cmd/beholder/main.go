package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var outPath string
var statusPath string

type supervisor struct {
	cmd       *exec.Cmd
	output    io.Writer
	status    io.Writer
	startedAt time.Time
	lastError error
}

func (sup *supervisor) Start() error {
	sup.startedAt = time.Now()
	sup.lastError = sup.cmd.Start()
	sup.WriteStatus()
	return sup.lastError
}

func (sup *supervisor) Wait() error {
	sup.lastError = sup.cmd.Wait()
	sup.WriteStatus()
	return sup.lastError
}

func (sup *supervisor) WriteStatus() {
	writeOk(sup.status.Write(snapshot(sup).Bytes()))
}

// status of a supervisor
type stat struct {
	Running    bool
	Error      string
	Exited     bool
	Pid        int
	Desc       string
	Success    bool
	Elapsed    time.Duration
	SystemTime time.Duration
	UserTime   time.Duration
}

func snapshot(sup *supervisor) *stat {
	result := &stat{
		Elapsed: time.Now().Sub(sup.startedAt),
	}

	if state := sup.cmd.ProcessState; state != nil {
		result.Exited = state.Exited()
		result.Pid = state.Pid()
		result.Desc = state.String()
		result.Success = state.Success()
		result.SystemTime = state.SystemTime()
		result.UserTime = state.UserTime()
	} else {
		if cmd := sup.cmd; cmd.Process != nil {
			result.Pid = cmd.Process.Pid
			result.Running = true
		}
	}

	if sup.lastError != nil {
		result.Error = fmt.Sprintf("%v", sup.lastError)
	}

	return result
}

func (s *stat) Bytes() []byte {
	data, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return append(data, '\n')
}

func (s *stat) String() string {
	return string(s.Bytes())
}

const usage = `Usage of %s:

  %s [OPTIONS] PROGRAM [ARG...]

Execute the given program, printing its output to either STDERR or a file.
Print the status of the execution to STDOUT or a file as newline-seperated JSON.

This utility is intended to make it easy to gather the status of a long-running
command once it has finished.

Examples:
  # preview the sort of output you'll get
  %s echo hello world

  # run make disowned, so you can retrieve the build status later
  nohup %s -status make.status -output make.log make

Flags:
`

func init() {
	flag.Usage = func() {
		name := os.Args[0]
		fmt.Fprintf(os.Stderr, usage, name, name, name, name)
		flag.PrintDefaults()
	}
	flag.StringVar(&outPath, "out", "", "Location to write out program output")
	flag.StringVar(&statusPath, "status", "", "Location to write out program exit status")
	flag.Parse()
}

func main() {
	var err error

	output := os.Stderr
	if outPath != "" {
		output, err = os.Create(outPath)
		if err != nil {
			log.Fatalf("Cannot open output file: %v", err)
		}
	}

	status := os.Stdout
	if statusPath != "" {
		status, err = os.Create(statusPath)
		if err != nil {
			log.Fatalf("Cannot open status file: %v", err)
		}
	}

	args := flag.Args()

	if len(args) == 0 {
	}

	name := args[0]
	rest := args[1:]

	cmd := exec.Command(name, rest...)
	cmd.Stdout = output
	cmd.Stderr = output

	sup := &supervisor{
		cmd:    cmd,
		output: output,
		status: status,
	}

	if err := sup.Start(); err != nil {
		log.Fatalln(err)
	}

	if err := sup.Wait(); err != nil {
		log.Fatalln(err)
	}
}

func writeOk(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
