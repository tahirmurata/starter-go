package manager

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Manager interface {
	// Add registers a function to be called during shutdown in FIFO order
	Add(stopFunc func())

	// WaitForStop calls Stop and waits until it is done
	WaitForStop()

	// WaitForSignal blocks until SIGINT or SIGTERM is received, then calls Stop
	WaitForSignal()
}

type manager struct {
	stopFuncs []func()
	mu        sync.Mutex
	done      chan struct{}
}

var instance *manager
var once sync.Once

func Init() Manager {
	once.Do(func() {
		instance = &manager{
			stopFuncs: []func(){},
			done:      make(chan struct{}),
		}
	})
	return instance
}

func (m *manager) Add(stopFunc func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopFuncs = append(m.stopFuncs, stopFunc)
}

func (m *manager) stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Call functions in FILO (Last-In-First-Out) order
	for i := len(m.stopFuncs) - 1; i >= 0; i-- {
		slog.Info("stopping service", "index", i)
		stopFunc := m.stopFuncs[i]
		stopFunc()
	}

	// Clear functions list after stopping
	m.stopFuncs = nil
	close(m.done)
}

func (m *manager) WaitForStop() {
	// Initiate graceful shutdown
	m.stop()

	// Wait for shutdown to complete
	<-m.done
	slog.Info("graceful shutdown complete")
}

func (m *manager) WaitForSignal() {
	// Create notification channel for signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-sigCh
	slog.Info("received signal, starting shutdown", "signal", sig)

	// Initiate graceful shutdown
	m.stop()

	// Wait for shutdown to complete
	<-m.done
	slog.Info("graceful shutdown complete")
}
