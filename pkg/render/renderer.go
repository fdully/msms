// Package render defines rendering functionality.
package render

import (
	"bytes"
	"sync"
)

// Renderer is a structure that knows how to perform safe HTTP rendering.
type Renderer struct {
	pool *sync.Pool
}

// NewRenderer returns an instantiated renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 1024))
			},
		},
	}
}
