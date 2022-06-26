package status

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) AddStatus(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) DelStatus(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
