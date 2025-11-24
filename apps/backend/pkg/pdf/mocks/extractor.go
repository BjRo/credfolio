package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type MockPDFExtractor struct {
	mock.Mock
}

func (m *MockPDFExtractor) ExtractText(reader io.Reader) (string, error) {
	args := m.Called(reader)
	return args.String(0), args.Error(1)
}
