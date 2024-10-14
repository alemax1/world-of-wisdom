package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"go.uber.org/zap"

	"github.com/alemax1/world-of-wisdom/internal/model"
)

var (
	ErrIvalidNonce      = errors.New("invalid nonce")
	ErrInvalidChallenge = errors.New("invalid challenge")
)

type Solver interface {
	SolveChallenge(difficulty uint8, challenge []byte) (uint64, error)
}

type Client struct {
	solver Solver
	addr   string
	log    *zap.Logger
}

func New(
	solver Solver,
	addr string,
	log *zap.Logger,
) *Client {
	return &Client{
		solver: solver,
		addr:   addr,
		log:    log,
	}
}

func (c Client) Run() error {
	conn, err := net.Dial(model.TCPName, c.addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	defer func() {
		if e := conn.Close(); e != nil {
			if err != nil {
				err = fmt.Errorf("%w, close conn: %w", err, e)

				return
			}

			err = fmt.Errorf("close conn: %w", e)
		}
	}()

	buf := make([]byte, model.MaxChallengeSize)
	if _, err := conn.Read(buf); err != nil {
		return fmt.Errorf("read max challenge size: %w", err)
	}

	challengeSize := binary.BigEndian.Uint32(buf[:model.MaxChallengeSize]) + uint32(model.MaxChallengeDifficultySize)

	challengeBuf := make([]byte, challengeSize)
	if _, err := conn.Read(challengeBuf); err != nil {
		return fmt.Errorf("read challenge: %w", err)
	}

	challengeDiff := uint8(challengeBuf[challengeSize-1])

	nonce, err := c.solver.SolveChallenge(challengeDiff, challengeBuf[:challengeSize-1])
	if err != nil {
		return fmt.Errorf("solve challenge: %w", err)
	}

	nonceBuf := make([]byte, model.MaxNonceSize)
	binary.BigEndian.PutUint64(nonceBuf, nonce)

	if _, err := conn.Write(nonceBuf); err != nil {
		return fmt.Errorf("write nonce: %w", err)
	}

	challenge, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("read challenge: %w", err)
	}

	if len(challenge) == 0 {
		return ErrInvalidChallenge
	}

	c.log.Info("resolve challenge", zap.ByteString("data", challenge))

	return nil
}
