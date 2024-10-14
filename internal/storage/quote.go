package storage

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
)

var ErrNoDataFound = errors.New("no data in provided file")

type Storage struct {
	FilePath string
	Quotes   []quote
}

type quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func New(filePath string) (*Storage, error) {
	storage := &Storage{
		FilePath: filePath,
	}

	if err := storage.fill(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) fill() error {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return fmt.Errorf("read from file: %w", err)
	}

	var quotes []quote

	if err := json.Unmarshal(data, &quotes); err != nil {
		return fmt.Errorf("unmarshal data from file: %w", err)
	}

	if len(quotes) == 0 {
		return ErrNoDataFound
	}

	s.Quotes = quotes

	return nil
}

func (s *Storage) GetRandQuoteBytes() ([]byte, error) {
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(s.Quotes)-1)))
	if err != nil {
		return nil, fmt.Errorf("generate rand index: %w", err)
	}

	quote := s.Quotes[randIndex.Int64()]

	data, err := json.Marshal(quote)
	if err != nil {
		return nil, fmt.Errorf("marhsal quote: %w", err)
	}

	return data, nil
}
