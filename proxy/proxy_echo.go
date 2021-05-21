package proxy

func Echo(ctx *ProxyContext, pes *ProxyEndpoints) error {
	for {
		select {
		case m := <-pes.Downstream.Rx():
			pes.Downstream.Tx() <- m
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
