import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import ProfileHeader from './ProfileHeader';

describe('ProfileHeader', () => {
  const mockProfile = {
    id: '123',
    full_name: 'John Doe',
    email: 'john@example.com',
  };

  it('renders user information', () => {
    render(<ProfileHeader profile={mockProfile} />);
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
    // Initials
    expect(screen.getByText('J')).toBeInTheDocument();
  });

  it('renders download CV button with correct link', () => {
    render(<ProfileHeader profile={mockProfile} />);
    const link = screen.getByRole('link', { name: /Download CV/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', expect.stringContaining('/api/profile/cv?user_id=123'));
  });

  it('renders nothing if profile is missing', () => {
    const { container } = render(<ProfileHeader profile={null} />);
    expect(container).toBeEmptyDOMElement();
  });
});

