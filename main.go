package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/borgstrom/ebgb/emulator"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <rom>", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read %s: %s", os.Args[1], err)
	}
	defer f.Close()

	ctx, cancel := ContextWithCancelAndSignals(context.Background())
	defer cancel()

	e := emulator.New(f)
	e.Run(ctx)
}

// ContextWithCancelAndSignals returns a context and a cancel function.  The context will
// be cancelled upon a SIGTERM or SIGINT being received by the current process.
//
// NOTE: This installs signal handlers for only the *first* SIGTERM or SIGINT
// received.  It cancels the returned context after receiving the signal and
// then removes the signal handlers so that successive signals are handled in
// the default way in case there is a bug in the program.
func ContextWithCancelAndSignals(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ready := make(chan struct{})

	go func() {
		// Create a channel for receiving the signals
		shutdown := make(chan os.Signal)

		// Bind the channel to the signals
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

		// Reset the signal handler when we're done
		defer signal.Reset(os.Interrupt, syscall.SIGTERM)

		// Close the ready channel indicating we are ready to handle the signal
		close(ready)

		select {
		case <-parent.Done():
			return
		case sig := <-shutdown:
			// If we receive a signal, log it and cancel the context
			log.Printf("Caught signal: %s", sig)
			cancel()
		}
	}()

	// Wait for the signal handler to be installed before returning
	<-ready
	return ctx, cancel
}
