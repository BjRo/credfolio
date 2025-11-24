import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import CredibilityHighlights from "./CredibilityHighlights";
import type { CredibilityHighlight } from "../../lib/api/generated/models/CredibilityHighlight";

describe("CredibilityHighlights", () => {
	it("renders nothing when highlights are empty", () => {
		const { container } = render(<CredibilityHighlights highlights={[]} />);
		expect(container.firstChild).toBeNull();
	});

	it("renders nothing when highlights are undefined", () => {
		const { container } = render(
			<CredibilityHighlights highlights={undefined} />,
		);
		expect(container.firstChild).toBeNull();
	});

	it("renders credibility highlights with positive sentiment", () => {
		const highlights: Array<CredibilityHighlight> = [
			{
				quote: "Excellent team player",
				sentiment: "POSITIVE" as const,
			},
		];
		render(<CredibilityHighlights highlights={highlights} />);
		expect(screen.getByText("Employer Feedback")).toBeInTheDocument();
		expect(screen.getByText(/Excellent team player/)).toBeInTheDocument();
	});

	it("renders credibility highlights with neutral sentiment", () => {
		const highlights: Array<CredibilityHighlight> = [
			{
				quote: "Reliable worker",
				sentiment: "NEUTRAL" as const,
			},
		];
		render(<CredibilityHighlights highlights={highlights} />);
		expect(screen.getByText("Employer Feedback")).toBeInTheDocument();
		expect(screen.getByText(/Reliable worker/)).toBeInTheDocument();
	});

	it("renders multiple highlights", () => {
		const highlights: Array<CredibilityHighlight> = [
			{
				quote: "First quote",
				sentiment: "POSITIVE" as const,
			},
			{
				quote: "Second quote",
				sentiment: "NEUTRAL" as const,
			},
		];
		render(<CredibilityHighlights highlights={highlights} />);
		expect(screen.getByText(/First quote/)).toBeInTheDocument();
		expect(screen.getByText(/Second quote/)).toBeInTheDocument();
	});
});
