import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import CredibilitySection from './CredibilitySection';

describe('CredibilitySection', () => {
  it('renders stats correctly', () => {
    const companies = [
      { roles: [{ is_verified: true }, { is_verified: true }] },
      { roles: [{ is_verified: false }] },
    ];
    // Total roles = 3, Verified = 2. Rate = 67%

    render(<CredibilitySection companies={companies} />);

    expect(screen.getByText('2')).toBeInTheDocument(); // Verified Count
    expect(screen.getByText('Verified Roles')).toBeInTheDocument();

    expect(screen.getByText('67%')).toBeInTheDocument(); // Rate
    expect(screen.getByText('Verification Rate')).toBeInTheDocument();
  });

  it('renders nothing if no verified roles', () => {
    const companies = [
        { roles: [{ is_verified: false }] }
    ];
    const { container } = render(<CredibilitySection companies={companies} />);
    expect(container).toBeEmptyDOMElement();
  });

  it('handles empty companies list safely', () => {
      const { container } = render(<CredibilitySection companies={[]} />);
      expect(container).toBeEmptyDOMElement();
  });
});

