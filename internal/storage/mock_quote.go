package storage

type StorageMock struct{}

func NewMock() (*StorageMock, error) {
	return &StorageMock{}, nil
}

func (s *StorageMock) GetRandQuoteBytes() ([]byte, error) {
	return []byte("Some quote"), nil
}
