package proxy

type Endpoint interface {
	Rx() <-chan *Msg
	Tx() chan<- *Msg
}

type ProxyEndpoints struct {
	Upstream Endpoint
	Downstream Endpoint
}

type DuplexEndpoint struct {
	Receive chan *Msg
	Transmit chan *Msg
}

func NewDuplexEndpoint(rxSize int, txSize int) *DuplexEndpoint {
	return &DuplexEndpoint{
		Receive: make(chan *Msg, rxSize),
		Transmit: make(chan *Msg, txSize),
	}
}

func (it *DuplexEndpoint) Rx() <-chan *Msg {
	return it.Receive
}

func (it *DuplexEndpoint) Tx() chan<- *Msg {
	return it.Transmit
}

func (it *DuplexEndpoint) Close() {
	close(it.Receive)
	close(it.Transmit)
}
