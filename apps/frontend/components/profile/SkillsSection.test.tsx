import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { SkillsSection } from "./SkillsSection";

describe("SkillsSection", () => {
	const mockSkills = [
		"Go",
		"TypeScript",
		"React",
		"PostgreSQL",
		"Docker",
		"Kubernetes",
	];

	// T092: Unit test for SkillsSection component when given skills displays aggregated list
	it("when given skills, displays aggregated list", () => {
		// Arrange & Act
		render(<SkillsSection skills={mockSkills} />);

		// Assert
		expect(screen.getByText("Go")).toBeInTheDocument();
		expect(screen.getByText("TypeScript")).toBeInTheDocument();
		expect(screen.getByText("React")).toBeInTheDocument();
		expect(screen.getByText("PostgreSQL")).toBeInTheDocument();
		expect(screen.getByText("Docker")).toBeInTheDocument();
		expect(screen.getByText("Kubernetes")).toBeInTheDocument();
	});

	it("when given empty skills, shows placeholder", () => {
		// Arrange & Act
		render(<SkillsSection skills={[]} />);

		// Assert
		expect(screen.getByText(/no skills added yet/i)).toBeInTheDocument();
	});

	it("displays skills as badges/tags", () => {
		// Arrange & Act
		render(<SkillsSection skills={mockSkills} />);

		// Assert
		const skillElements = screen.getAllByRole("listitem");
		expect(skillElements.length).toBe(6);
	});

	it("removes duplicate skills from display", () => {
		// Arrange
		const skillsWithDuplicates = [
			"Go",
			"TypeScript",
			"Go",
			"React",
			"TypeScript",
		];

		// Act
		render(<SkillsSection skills={skillsWithDuplicates} />);

		// Assert - Should only show unique skills
		const goElements = screen.getAllByText("Go");
		expect(goElements.length).toBe(1);

		const tsElements = screen.getAllByText("TypeScript");
		expect(tsElements.length).toBe(1);
	});

	it("displays skill count", () => {
		// Arrange & Act
		render(<SkillsSection skills={mockSkills} />);

		// Assert
		expect(screen.getByText(/6 skills/i)).toBeInTheDocument();
	});

	it("sorts skills alphabetically", () => {
		// Arrange
		const unsortedSkills = ["Zebra", "Apple", "Mango"];

		// Act
		render(<SkillsSection skills={unsortedSkills} />);

		// Assert
		const skillElements = screen.getAllByRole("listitem");
		expect(skillElements[0]).toHaveTextContent("Apple");
		expect(skillElements[1]).toHaveTextContent("Mango");
		expect(skillElements[2]).toHaveTextContent("Zebra");
	});
});
