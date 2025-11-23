import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import FileUploader from './FileUploader';

// Mock the API module
vi.mock('../../services/api', () => ({
  uploadReferenceLetter: vi.fn(),
}));

import { uploadReferenceLetter } from '../../services/api';

describe('FileUploader', () => {
  it('renders upload area', () => {
    render(<FileUploader />);
    expect(screen.getByText(/Drag & Drop or Click to Upload/i)).toBeInTheDocument();
    expect(screen.getByText(/PDF files only/i)).toBeInTheDocument();
  });

  it('shows error for non-PDF files', async () => {
    render(<FileUploader />);
    const input = screen.getByLabelText(/Drag & Drop or Click to Upload/i);

    const file = new File(['dummy content'], 'test.txt', { type: 'text/plain' });
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText(/Only PDF files are supported/i)).toBeInTheDocument();
    });
  });

  it('handles successful upload', async () => {
    const mockUpload = vi.mocked(uploadReferenceLetter);
    mockUpload.mockResolvedValueOnce({ status: 'processing', user_id: '123' });

    render(<FileUploader />);
    const input = screen.getByLabelText(/Drag & Drop or Click to Upload/i);

    const file = new File(['dummy content'], 'test.pdf', { type: 'application/pdf' });
    fireEvent.change(input, { target: { files: [file] } });

    expect(screen.getByText(/Uploading.../i)).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText(/File uploaded successfully! Processing.../i)).toBeInTheDocument();
    });
  });

  it('handles upload failure', async () => {
    const mockUpload = vi.mocked(uploadReferenceLetter);
    mockUpload.mockRejectedValueOnce(new Error('Network error'));

    render(<FileUploader />);
    const input = screen.getByLabelText(/Drag & Drop or Click to Upload/i);

    const file = new File(['dummy content'], 'test.pdf', { type: 'application/pdf' });
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText(/Upload failed. Please try again./i)).toBeInTheDocument();
    });
  });
});

