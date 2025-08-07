package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransport(listenAddr)

	// ✅ Fix: expected, actual order
	assert.Equal(t, listenAddr, tr.listenAddress)

	// ✅ Run ListenAndAccept and ensure it doesn't return an error
	go func() {
		err := tr.ListenAndAccept()
		assert.Nil(t, err)
	}()

	// ✅ Short delay to let goroutine start
	// In real code you'd use sync.WaitGroup or context
	select {}
}
