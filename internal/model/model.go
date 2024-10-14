package model

const (
	TCPName                           = "tcp"
	MaxChallengeSize           uint32 = 4 // Maximum length of challenge size in bytes.
	MaxNonceSize               uint64 = 8 // Maximum length of nonce size in bytes.
	MaxChallengeDifficultySize uint8  = 1 // Maximum length of challenge difficulty in byte.
)
