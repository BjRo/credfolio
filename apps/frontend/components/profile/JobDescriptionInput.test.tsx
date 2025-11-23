import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import JobDescriptionInput from './JobDescriptionInput';

// Mock API
vi.mock('../../services/api', () => ({
  tailorProfile: vi.fn(),
}));

import { tailorProfile } from '../../services/api';

describe('JobDescriptionInput', () => {
  const mockOnTailored = vi.fn();
  const userID = '123';

  it('renders input area and button', () => {
    render(<JobDescriptionInput userID={userID} onTailored={mockOnTailored} />);
    expect(screen.getByPlaceholderText(/Paste job description here/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Tailor Profile/i })).toBeInTheDocument();
  });

  it('disables button when input is empty', () => {
    render(<JobDescriptionInput userID={userID} onTailored={mockOnTailored} />);
    const button = screen.getByRole('button', { name: /Tailor Profile/i });
    expect(button).toBeDisabled();
  });

  it('calls api on submit', async () => {
    const mockTailor = vi.mocked(tailorProfile);
    mockTailor.mockResolvedValueOnce({ match_score: 90 });

    render(<JobDescriptionInput userID={userID} onTailored={mockOnTailored} />);

    const textarea = screen.getByPlaceholderText(/Paste job description here/i);
    fireEvent.change(textarea, { target: { value: 'Looking for React dev' } });

    const button = screen.getByRole('button', { name: /Tailor Profile/i });
    expect(button).not.toBeDisabled();

    fireEvent.click(button);

    expect(screen.getByText(/Analyzing.../i)).toBeInTheDocument();

    await waitFor(() => {
      expect(mockTailor).toHaveBeenCalledWith(userID, 'Looking for React dev');
      expect(mockOnTailored).toHaveBeenCalledWith({ match_score: 90 });
      expect(screen.queryByText(/Analyzing.../i)).not.toBeInTheDocument();
    });
  });

  it('handles api error', async () => {
    const mockTailor = vi.mocked(tailorProfile);
    mockTailor.mockRejectedValueOnce(new Error('API Error'));

    render(<JobDescriptionInput userID={userID} onTailored={mockOnTailored} />);

    const textarea = screen.getByPlaceholderText(/Paste job description here/i);
    fireEvent.change(textarea, { target: { value: 'Looking for React dev' } });
    fireEvent.click(screen.getByRole('button', { name: /Tailor Profile/i }));

    await waitFor(() => {
      expect(screen.getByText(/Failed to tailor profile/i)).toBeInTheDocument();
    });
  });
});

