package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRandChallengeUniqueness(t *testing.T) {
	m := Manager{ChallengeSize: 32}

	c1, err := m.CreateRandChallenge()
	assert.NoError(t, err)
	c2, err := m.CreateRandChallenge()
	assert.NoError(t, err)
	c3, err := m.CreateRandChallenge()
	assert.NoError(t, err)

	assert.NotEqual(t, c1, c2)
	assert.NotEqual(t, c1, c3)
	assert.NotEqual(t, c2, c3)
}

func TestCreateRandChallengeLengthEquals(t *testing.T) {
	var challengeSize uint32 = 0
	m := Manager{ChallengeSize: challengeSize}

	c1, err := m.CreateRandChallenge()
	assert.NoError(t, err)
	c2, err := m.CreateRandChallenge()
	assert.NoError(t, err)

	assert.Len(t, c1, int(challengeSize))
	assert.Len(t, c2, int(challengeSize))

	challengeSize = 64
	m.ChallengeSize = 64

	c3, err := m.CreateRandChallenge()
	assert.NoError(t, err)
	c4, err := m.CreateRandChallenge()
	assert.NoError(t, err)

	assert.Len(t, c3, int(challengeSize))
	assert.Len(t, c4, int(challengeSize))
}

func TestSolveAndValidateChallenge(t *testing.T) {
	var diff uint8 = 16
	challenge := []byte("asdYTtWfsha123ASd")
	p := Solver{}
	nonce, err := p.SolveChallenge(diff, challenge)
	assert.NoError(t, err, "soolve challenge")

	challengeWithNonce := binary.BigEndian.AppendUint64(challenge, nonce)

	hashSum := sha256.Sum256(challengeWithNonce)

	log.Println(hashSum)

	m := Manager{}
	assert.True(t, m.ValidateChallenge(diff, nonce, challenge))
	assert.False(t, m.ValidateChallenge(diff+1, nonce, challenge))
}

func TestCountZeros(t *testing.T) {
	hash := [32]byte{}
	result := countZeros(hash)

	assert.Equal(t, 255, int(result))
}
