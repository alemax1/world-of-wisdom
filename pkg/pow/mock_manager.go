package pow

type MockManager struct {
}

func NewMockManager() *MockManager {
	return &MockManager{}
}

func (m MockManager) CreateRandChallenge() ([]byte, error) {
	return []byte("HELLO"), nil
}

func (m MockManager) ValidateChallenge(_ uint8, _ uint64, _ []byte) bool {
	return true
}
