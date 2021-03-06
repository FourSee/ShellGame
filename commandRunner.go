package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type commandRunner struct {
	command           string
	args              []string
	startNotification bool
}

func (cr *commandRunner) run() {
	var (
		cmd *exec.Cmd
		err error
	)

	argList := strings.Join(cr.args, ",")

	cmd = exec.Command(cr.command, argList)

	setupPrinter(cmd)

	start := time.Now()

	if err = cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}

	defer timeRan(cmd.Process.Pid, start)

	ss := StartShell{Pid: cmd.Process.Pid, SendNotification: cr.startNotification, StartDate: start}

	fmt.Println("starting!", ss)

	if err := cmd.Wait(); err != nil {
		fmt.Println("Waiting")

		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
}

func setupPrinter(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatalln("Unable to get output from program")
	}

	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
}

func timeRan(pid int, st time.Time) {
	t := time.Now()

	es := EndShell{Pid: pid}
	es.setElapsed(st, t)

	fmt.Println("End Shell:", es)
	fmt.Println("Time Elapsed:", es.Elapsed)
}
