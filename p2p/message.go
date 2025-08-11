package p2p

const (
	IncomingMessage = 0x2
	IncomingStream  = 0x1
)

// RPC holds any arbitrary data that is being sent over the
// each transport between two nodes in the network
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
