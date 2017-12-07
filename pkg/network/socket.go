package network

func NewSocket(config Config) listenInterface {
	switch config.Proto {
	case Tcp:
		return &tcpSocket{
			Config: config,
		}
	default:
		return &tcpSocket{
			Config: config,
		}
	}

}
