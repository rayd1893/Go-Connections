package render

import (
	"bytes"
	"sync"
)

type Renderer struct {
	pool *sync.Pool
}

func NewRenderer() *Renderer {
	return &Renderer{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 1024))
			},
		},
	}
}
