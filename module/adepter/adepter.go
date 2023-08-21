package adepter

import (
	"google.golang.org/grpc"

	"github.com/meowalien/RabbitGather-interest-crawler.git/fremwork/connect"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module"
)

// Adepter is the transfer request from connection to module.
type Adepter interface {
	// GRPC transfer the grpc call to module.
	GRPC(svc *grpc.Server)
	// RabbitMQ transfer the rabbitmq message to module.
	RabbitMQ(cp connect.ChannelPool)
}

type Constructor struct {
	Module module.Module
}

func (a Constructor) New() Adepter {
	return &adepter{
		md: a.Module,
	}
}

type adepter struct {
	md module.Module
}
