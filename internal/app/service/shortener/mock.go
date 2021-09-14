package shortener

import "github.com/stretchr/testify/mock"

type StoreMock struct {
	mock.Mock
}

var _ Store = (*StoreMock)(nil)

func (m *StoreMock) WriteURL(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *StoreMock) ReadURL(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}
