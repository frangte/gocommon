package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewOSExitWaiter creates a waiter that responds to OS shutdown signals (SIGINT, SIGTERM).
// It returns two channels:
//   - quiting: signals that shutdown has begun
//   - quited: signals that shutdown has completed
//
// The timeout parameter determines how long to wait before forcing shutdown.
func NewOSExitWaiter(ctx context.Context, timeout time.Duration) (quiting chan struct{}, quited chan struct{}) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	quiting = make(chan struct{}, 1)
	quited = make(chan struct{}, 1)
	go func() {
		<-quit
		quiting <- struct{}{}
		c, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		<-c.Done()
		quited <- struct{}{}
		close(quiting)
		close(quited)
	}()

	return
}
