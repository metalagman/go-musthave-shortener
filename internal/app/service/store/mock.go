package store

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

var _ Store = (*Mock)(nil)

func (m *Mock) WriteURL(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *Mock) ReadURL(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}
