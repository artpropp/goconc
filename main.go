package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
)

const cmdDelimiter = "::"

type cmdArgs struct {
	name string
	args []string
}

func parseArgs(args []string) []cmdArgs {
	if len(args) == 0 {
		return nil
	}
	var cmds []cmdArgs
	var cmd cmdArgs
	for _, arg := range args {
		switch {
		case arg == cmdDelimiter:
			cmds = append(cmds, cmd)
			cmd = cmdArgs{}
		case cmd.name == "":
			cmd.name = arg
		default:
			cmd.args = append(cmd.args, arg)
		}
	}
	cmds = append(cmds, cmd)
	return cmds
}

func runCmds(ctx context.Context, cmds []cmdArgs) {
	wg := &sync.WaitGroup{}
	for i, args := range cmds {
		wg.Add(1)
		go func(nr int, args cmdArgs) {
			log.Printf("Starting cmd %d: %s %s\n", nr, args.name, strings.Join(args.args, " "))
			cmd := exec.CommandContext(ctx, args.name, args.args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			log.Printf("cmd %d is ready\n", nr)

			if err := cmd.Run(); err != nil {
				log.Printf("cmd %d with error = %s\n", nr, err)
			}
			wg.Done()
		}(i+1, args)
	}
	wg.Wait()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Canceling with Ctrl+C")
		cancel()
	}()

	cmds := parseArgs(os.Args[1:])
	runCmds(ctx, cmds)
}
