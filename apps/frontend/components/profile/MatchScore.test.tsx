import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { MatchScore } from "./MatchScore";

describe("MatchScore", () => {
	// T113: Unit test for MatchScore component when given score displays percentage
	it("when given score, displays percentage", () => {
		// Arrange & Act
		render(<MatchScore score={0.85} />);

		// Assert
		expect(screen.getByText("85%")).toBeInTheDocument();
	});

	it("displays high match indicator for scores >= 0.7", () => {
		// Arrange & Act
		render(<MatchScore score={0.75} />);

		// Assert
		expect(screen.getByText(/strong match/i)).toBeInTheDocument();
	});

	it("displays medium match indicator for scores 0.4-0.7", () => {
		// Arrange & Act
		render(<MatchScore score={0.55} />);

		// Assert
		expect(screen.getByText(/moderate match/i)).toBeInTheDocument();
	});

	it("displays low match indicator for scores < 0.4", () => {
		// Arrange & Act
		render(<MatchScore score={0.25} />);

		// Assert
		expect(screen.getByText(/limited match/i)).toBeInTheDocument();
	});

	it("displays summary when provided", () => {
		// Arrange & Act
		render(
			<MatchScore score={0.85} summary="Strong match based on Go experience" />,
		);

		// Assert
		expect(
			screen.getByText(/strong match based on go experience/i),
		).toBeInTheDocument();
	});

	it("rounds percentage correctly", () => {
		// Arrange & Act
		render(<MatchScore score={0.867} />);

		// Assert
		expect(screen.getByText("87%")).toBeInTheDocument();
	});
});
