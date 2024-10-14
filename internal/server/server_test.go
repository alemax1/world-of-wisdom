package server

import (
	"net"
	"testing"

	"github.com/alemax1/world-of-wisdom/internal/config"
	"github.com/alemax1/world-of-wisdom/internal/storage"
	"github.com/alemax1/world-of-wisdom/pkg/logger"
	"github.com/alemax1/world-of-wisdom/pkg/pow"

	"github.com/stretchr/testify/assert"
)

func TestHandleConnection(t *testing.T) {
	m := pow.NewMockManager()
	st, _ := storage.NewMock()
	s := New(m, st, *config.New(), logger.NewMockLogger())

	l, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	clConn, err := net.Dial("tcp", l.Addr().String())
	assert.NoError(t, err)

	conn, err := l.Accept()
	assert.NoError(t, err)

	go s.handle(conn)

	buf := make([]byte, 128)
	_, err = clConn.Read(buf)
	assert.NoError(t, err)

	_, err = clConn.Write([]byte{1, 2, 3, 4})
	assert.NoError(t, err)

	buf = make([]byte, 10)
	_, err = clConn.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, "Some quote", string(buf)) // Some quote returns from MockStorage
}
