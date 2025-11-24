package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFileType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		filename string
		wantErr  bool
		errCode  int
	}{
		{
			name:     "valid PDF with application/pdf MIME type",
			mimeType: "application/pdf",
			filename: "test.pdf",
			wantErr:  false,
		},
		{
			name:     "valid PDF with application/x-pdf MIME type",
			mimeType: "application/x-pdf",
			filename: "test.pdf",
			wantErr:  false,
		},
		{
			name:     "valid PDF extension with missing MIME type",
			mimeType: "",
			filename: "test.pdf",
			wantErr:  false,
		},
		{
			name:     "valid PDF extension with invalid MIME type",
			mimeType: "text/plain",
			filename: "test.pdf",
			wantErr:  false,
		},
		{
			name:     "valid MIME type with wrong extension",
			mimeType: "application/pdf",
			filename: "test.txt",
			wantErr:  false,
		},
		{
			name:     "invalid MIME type and extension",
			mimeType: "text/plain",
			filename: "test.txt",
			wantErr:  true,
			errCode:  ErrorCodeInvalidFileType,
		},
		{
			name:     "case insensitive MIME type",
			mimeType: "APPLICATION/PDF",
			filename: "test.pdf",
			wantErr:  false,
		},
		{
			name:     "case insensitive extension",
			mimeType: "application/pdf",
			filename: "test.PDF",
			wantErr:  false,
		},
		{
			name:     "MIME type with whitespace",
			mimeType: "  application/pdf  ",
			filename: "test.pdf",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileType(tt.mimeType, tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
				if valErr, ok := err.(*ValidationError); ok {
					assert.Equal(t, tt.errCode, valErr.ErrorCode)
					assert.Equal(t, "file", valErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name        string
		fileSize    int64
		maxSize     int64
		wantErr     bool
		errCode     int
		containsMsg string
	}{
		{
			name:     "valid file size within limit",
			fileSize: 5 * 1024 * 1024,
			maxSize:  10 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "valid file size at limit",
			fileSize: 10 * 1024 * 1024,
			maxSize:  10 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "valid small file",
			fileSize: 1024,
			maxSize:  10 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:        "file size exceeds limit",
			fileSize:    15 * 1024 * 1024,
			maxSize:     10 * 1024 * 1024,
			wantErr:     true,
			errCode:     ErrorCodeFileTooLarge,
			containsMsg: "exceeds maximum allowed size",
		},
		{
			name:        "empty file",
			fileSize:    0,
			maxSize:     10 * 1024 * 1024,
			wantErr:     true,
			errCode:     ErrorCodeInvalidRequest,
			containsMsg: "file is empty",
		},
		{
			name:        "negative file size",
			fileSize:    -1,
			maxSize:     10 * 1024 * 1024,
			wantErr:     true,
			errCode:     ErrorCodeInvalidRequest,
			containsMsg: "file is empty",
		},
		{
			name:        "file size just over limit",
			fileSize:    10*1024*1024 + 1,
			maxSize:     10 * 1024 * 1024,
			wantErr:     true,
			errCode:     ErrorCodeFileTooLarge,
			containsMsg: "exceeds maximum allowed size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileSize(tt.fileSize, tt.maxSize)
			if tt.wantErr {
				assert.Error(t, err)
				if valErr, ok := err.(*ValidationError); ok {
					assert.Equal(t, tt.errCode, valErr.ErrorCode)
					assert.Equal(t, "file", valErr.Field)
					if tt.containsMsg != "" {
						assert.Contains(t, valErr.Message, tt.containsMsg)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
