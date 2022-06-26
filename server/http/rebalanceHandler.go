package http

import (
	"bytes"
	"net/http"
)

type rebalanceHandler struct {
	*Server
}

func (h *rebalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	go h.rebalance()
}

func (h *rebalanceHandler) rebalance() {
	s := h.NewScanner()
	defer s.Close()

	c := &http.Client{}
	for s.Scan() {
		key := s.Key()
		address, ok := h.ShouldProcess(key)
		if !ok {
			request, _ := http.NewRequest(http.MethodPut, "https://"+address+":12345/cache/"+key, bytes.NewReader(s.Value()))
			c.Do(request)
			h.Del(key)
		}
	}
}
