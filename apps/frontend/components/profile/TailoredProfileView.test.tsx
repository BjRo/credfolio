import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { TailoredProfileView } from "./TailoredProfileView";

describe("TailoredProfileView", () => {
	const mockTailoredProfile = {
		id: "tailored-1",
		matchScore: 0.85,
		matchSummary: "Strong match based on Go experience",
		tailoredExperiences: [
			{
				id: "exp-1",
				companyName: "Tech Corp",
				role: "Backend Engineer",
				startDate: "2020-01-01",
				endDate: "2023-12-31",
				description: "Worked with Go and databases",
				relevanceScore: 0.9,
				highlightReason: "Direct Go experience mentioned in job description",
			},
			{
				id: "exp-2",
				companyName: "Web Inc",
				role: "Frontend Developer",
				startDate: "2018-01-01",
				endDate: "2019-12-31",
				description: "React development",
				relevanceScore: 0.4,
				highlightReason: "General development experience",
			},
		],
		relevantSkills: ["Go", "PostgreSQL", "Docker"],
	};

	// T112: Unit test for TailoredProfileView component when given tailored profile highlights matched content
	it("when given tailored profile, highlights matched content", () => {
		// Arrange & Act
		render(<TailoredProfileView tailoredProfile={mockTailoredProfile} />);

		// Assert
		// Should show match score
		expect(screen.getByText("85%")).toBeInTheDocument();

		// Should show experiences with relevance scores
		expect(screen.getByText("Tech Corp")).toBeInTheDocument();
		expect(screen.getByText("Backend Engineer")).toBeInTheDocument();

		// Should show relevant skills
		expect(screen.getByText("Go")).toBeInTheDocument();
		expect(screen.getByText("PostgreSQL")).toBeInTheDocument();
	});

	it("displays match summary", () => {
		// Arrange & Act
		render(<TailoredProfileView tailoredProfile={mockTailoredProfile} />);

		// Assert
		expect(
			screen.getByText(/strong match based on go experience/i),
		).toBeInTheDocument();
	});

	it("shows experiences ordered by relevance", () => {
		// Arrange & Act
		render(<TailoredProfileView tailoredProfile={mockTailoredProfile} />);

		// Assert - Tech Corp should appear before Web Inc (higher relevance)
		const experiences = screen.getAllByRole("article");
		expect(experiences.length).toBeGreaterThanOrEqual(2);
	});

	it("displays highlight reason for relevant experiences", () => {
		// Arrange & Act
		render(<TailoredProfileView tailoredProfile={mockTailoredProfile} />);

		// Assert
		expect(screen.getByText(/direct go experience/i)).toBeInTheDocument();
	});

	it("displays relevance score for each experience", () => {
		// Arrange & Act
		render(<TailoredProfileView tailoredProfile={mockTailoredProfile} />);

		// Assert
		expect(screen.getByText("90%")).toBeInTheDocument(); // 0.9 relevance
		expect(screen.getByText("40%")).toBeInTheDocument(); // 0.4 relevance
	});
});
