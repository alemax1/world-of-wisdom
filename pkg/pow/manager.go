package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

var ErrReachMaxNonceSize = errors.New("reach maximum nonce size")

type Manager struct {
	ChallengeSize uint32
}

func NewManager(challengeSize uint32) *Manager {
	return &Manager{
		ChallengeSize: challengeSize,
	}
}

func (m Manager) CreateRandChallenge() ([]byte, error) {
	buf := make([]byte, m.ChallengeSize)

	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (m Manager) ValidateChallenge(difficulty uint8, nonce uint64, challenge []byte) bool {
	challengeWithNonce := binary.BigEndian.AppendUint64(challenge, nonce)

	hashSum := sha256.Sum256(challengeWithNonce)

	return countZeros(hashSum) >= uint16(difficulty)
}
