package http

import (
	"distribute_cache/cache/cache"
	"distribute_cache/server/cluster"
	"net/http"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.Handle("/cluster", s.clusterHandler())
	http.Handle("/rebalance", s.rebalanceHandler())
	http.ListenAndServe(s.Addr()+":12345", nil)
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{
		c,
		n,
	}
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}

func (s *Server) clusterHandler() http.Handler {
	return &clusterHandler{s}
}

func (s *Server) rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}
