package log

type Entry struct {
	Index int64 `json:"index,omitempty"`
	Term  int64 `json:"term,omitempty"`
	Data  KV    `json:"data"`
}

type KV struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
