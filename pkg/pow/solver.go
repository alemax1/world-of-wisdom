package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
)

type Solver struct{}

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) SolveChallenge(difficulty uint8, challenge []byte) (uint64, error) {
	var nonce uint64

	for {
		if nonce == math.MaxUint64 {
			return 0, ErrReachMaxNonceSize
		}

		challengeWithNonce := binary.BigEndian.AppendUint64(challenge, nonce)

		hashSum := sha256.Sum256(challengeWithNonce)

		if countZeros(hashSum) >= uint16(difficulty) {
			return nonce, nil
		}

		nonce++
	}
}
