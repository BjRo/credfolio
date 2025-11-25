import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { WorkExperienceCard } from "./WorkExperienceCard";

describe("WorkExperienceCard", () => {
	const mockExperience = {
		id: "exp-1",
		companyName: "Tech Corp",
		role: "Senior Software Engineer",
		startDate: "2020-01-15",
		endDate: "2023-12-31",
		description:
			"Led development of microservices architecture. Mentored junior engineers and improved deployment pipelines.",
		credibilityHighlights: [
			{
				quote: "Outstanding technical leadership and mentorship",
				sentiment: "POSITIVE" as const,
			},
			{
				quote: "Consistently exceeded expectations",
				sentiment: "POSITIVE" as const,
			},
		],
	};

	// T091: Unit test for WorkExperienceCard component when given experience displays credibility highlights
	it("when given experience with credibility highlights, displays them", () => {
		// Arrange & Act
		render(<WorkExperienceCard experience={mockExperience} />);

		// Assert
		expect(
			screen.getByText(/outstanding technical leadership/i),
		).toBeInTheDocument();
		expect(
			screen.getByText(/consistently exceeded expectations/i),
		).toBeInTheDocument();
	});

	it("displays company name and role", () => {
		// Arrange & Act
		render(<WorkExperienceCard experience={mockExperience} />);

		// Assert
		expect(screen.getByText("Tech Corp")).toBeInTheDocument();
		expect(screen.getByText("Senior Software Engineer")).toBeInTheDocument();
	});

	it("displays formatted date range", () => {
		// Arrange & Act
		render(<WorkExperienceCard experience={mockExperience} />);

		// Assert
		// Should display dates in readable format
		expect(screen.getByText(/jan 2020/i)).toBeInTheDocument();
		expect(screen.getByText(/dec 2023/i)).toBeInTheDocument();
	});

	it("displays description", () => {
		// Arrange & Act
		render(<WorkExperienceCard experience={mockExperience} />);

		// Assert
		expect(
			screen.getByText(/led development of microservices/i),
		).toBeInTheDocument();
	});

	it("when end date is null, shows Present", () => {
		// Arrange
		const currentExperience = {
			...mockExperience,
			endDate: undefined,
		};

		// Act
		render(<WorkExperienceCard experience={currentExperience} />);

		// Assert
		expect(screen.getByText(/present/i)).toBeInTheDocument();
	});

	it("when no credibility highlights, does not show highlights section", () => {
		// Arrange
		const experienceNoHighlights = {
			...mockExperience,
			credibilityHighlights: [],
		};

		// Act
		render(<WorkExperienceCard experience={experienceNoHighlights} />);

		// Assert
		expect(
			screen.queryByText(/credibility highlights/i),
		).not.toBeInTheDocument();
	});

	it("displays credibility highlights count badge", () => {
		// Arrange & Act
		render(<WorkExperienceCard experience={mockExperience} />);

		// Assert
		// Should show a count of how many highlights
		expect(screen.getByText(/2 endorsements/i)).toBeInTheDocument();
	});
});
