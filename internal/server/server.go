package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/alemax1/world-of-wisdom/internal/config"
	"github.com/alemax1/world-of-wisdom/internal/model"
)

type Challenger interface {
	CreateRandChallenge() ([]byte, error)
	ValidateChallenge(difficulty uint8, nonce uint64, challenge []byte) bool
}

type Logger interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
}

type Storage interface {
	GetRandQuoteBytes() ([]byte, error)
}

type Server struct {
	challenger Challenger
	storage    Storage
	config     config.Cfg
	log        Logger
	workers    chan struct{}
}

func New(
	challenger Challenger,
	storage Storage,
	cfg config.Cfg,
	log Logger,
) *Server {
	return &Server{
		challenger: challenger,
		storage:    storage,
		config:     cfg,
		log:        log,
		workers:    make(chan struct{}, cfg.WorkersCount),
	}
}

func (s Server) Run(ctx context.Context) (err error) {
	s.log.Info("server started", zap.Int("port", int(s.config.Port)))

	lis, err := net.Listen(model.TCPName, fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("listen port: %w", err)
	}

	defer func() {
		if e := lis.Close(); e != nil {
			if err != nil {
				err = fmt.Errorf("%w, close conn: %w", err, e)

				return
			}

			err = fmt.Errorf("close conn: %w", e)
		}

		close(s.workers)
	}()

	connPool := make(chan net.Conn, s.config.ConnectionsLimit)

	go func() {
		defer close(connPool)

		for {
			conn, err := lis.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.log.Error("accept connection", zap.Error(err))

				continue
			}

			if err := conn.SetDeadline(time.Now().Add(s.config.ConnectionTimeout)); err != nil {
				s.log.Error("set deadline", zap.Error(err))

				continue
			}

			connPool <- conn
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("server stopped")

			return nil

		case conn, ok := <-connPool:
			if !ok {
				return nil
			}

			s.workers <- struct{}{}
			go func() {
				s.handle(conn)
				<-s.workers
			}()
		}
	}
}

func (s Server) handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Error("close connection", zap.Error(err))
		}
	}()

	challenge, err := s.challenger.CreateRandChallenge()
	if err != nil {
		s.log.Error("create rand challenge", zap.Error(err))

		return
	}

	challengeDiff := s.config.ChallengeDifficulty
	challengeSize := uint32(len(challenge))
	bufferSize := model.MaxChallengeSize + challengeSize + uint32(model.MaxChallengeDifficultySize)

	buf := make([]byte, bufferSize)
	binary.BigEndian.PutUint32(buf[:model.MaxChallengeSize], challengeSize)

	offset := model.MaxChallengeSize + challengeSize
	buf[offset] = challengeDiff
	copy(buf[model.MaxChallengeSize:offset], challenge)

	if _, err = conn.Write(buf); err != nil {
		s.log.Error("write to conn", zap.Error(err))

		return
	}

	nonceBuf := make([]byte, model.MaxNonceSize)
	if _, err = conn.Read(nonceBuf); err != nil {
		fmt.Println("ERRRRRRRRRR")
		s.log.Error("read from conn", zap.Error(err))

		return
	}

	nonceInt := binary.BigEndian.Uint64(nonceBuf)

	if s.challenger.ValidateChallenge(challengeDiff, nonceInt, challenge) {
		quoteBytes, err := s.storage.GetRandQuoteBytes()
		if err != nil {
			s.log.Error("get quote", zap.Error(err))

			return
		}

		if _, err := conn.Write(quoteBytes); err != nil {
			s.log.Error("write quote", zap.Error(err))

			return
		}

		return
	}

	s.log.Warn("challenge wasn't validated",
		zap.Any("nonce", nonceInt),
		zap.Any("challenge diff", challengeDiff),
		zap.Any("challenge", challenge),
	)
}
