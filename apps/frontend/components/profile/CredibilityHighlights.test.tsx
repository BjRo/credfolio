import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { CredibilityHighlights } from "./CredibilityHighlights";

describe("CredibilityHighlights", () => {
	const mockHighlights = [
		{
			quote: "Outstanding technical leadership and mentorship",
			sentiment: "POSITIVE" as const,
		},
		{
			quote: "Consistently delivered high-quality work on time",
			sentiment: "POSITIVE" as const,
		},
		{
			quote: "Good team collaborator",
			sentiment: "NEUTRAL" as const,
		},
	];

	// T090: Unit test for CredibilityHighlights component when given highlights displays quotes
	it("when given highlights, displays quotes", () => {
		// Arrange & Act
		render(<CredibilityHighlights highlights={mockHighlights} />);

		// Assert
		expect(
			screen.getByText(/outstanding technical leadership/i),
		).toBeInTheDocument();
		expect(
			screen.getByText(/consistently delivered high-quality work/i),
		).toBeInTheDocument();
		expect(screen.getByText(/good team collaborator/i)).toBeInTheDocument();
	});

	it("when given positive highlights, displays with positive styling", () => {
		// Arrange
		const positiveHighlight = [
			{
				quote: "Exceptional performance",
				sentiment: "POSITIVE" as const,
			},
		];

		// Act
		render(<CredibilityHighlights highlights={positiveHighlight} />);

		// Assert
		const quoteElement = screen.getByText(/exceptional performance/i);
		expect(quoteElement).toBeInTheDocument();
		// Check for positive indicator (could be icon, color, etc.)
		expect(quoteElement.closest("[data-sentiment]")).toHaveAttribute(
			"data-sentiment",
			"POSITIVE",
		);
	});

	it("when given empty highlights, shows placeholder", () => {
		// Arrange & Act
		render(<CredibilityHighlights highlights={[]} />);

		// Assert
		expect(
			screen.getByText(/no credibility highlights yet/i),
		).toBeInTheDocument();
	});

	it("displays quote marks around quotes", () => {
		// Arrange
		const highlight = [
			{
				quote: "Great engineer",
				sentiment: "POSITIVE" as const,
			},
		];

		// Act
		render(<CredibilityHighlights highlights={highlight} />);

		// Assert
		// Should have quote marks or be styled as a blockquote
		expect(screen.getByRole("blockquote")).toBeInTheDocument();
	});
});
