package tcp

import (
	"distribute_cache/cache/cache"
	"distribute_cache/server/cluster"
	"net"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go s.process(conn)
	}
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{
		c,
		n,
	}
}
