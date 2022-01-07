package goq

import (
	"testing"
	"runtime"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	baseGo := 2
	assert.Equal(t, baseGo, runtime.NumGoroutine())
	
	q := New[int, int]()
	assert.Equal(t, baseGo + 1, runtime.NumGoroutine())

	<-q.Stop()
	assert.Equal(t, baseGo + 1, runtime.NumGoroutine())
}
