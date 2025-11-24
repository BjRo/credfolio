import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import WorkExperienceCard from "./WorkExperienceCard";
import type { WorkExperience } from "../../lib/api/generated/models/WorkExperience";

describe("WorkExperienceCard", () => {
	it("renders work experience with all fields", () => {
		const experience: WorkExperience = {
			id: "123",
			companyName: "Acme Corp",
			role: "Software Engineer",
			startDate: "2020-01-01",
			endDate: "2022-12-31",
			description: "Built awesome features",
		};
		render(<WorkExperienceCard experience={experience} />);
		expect(screen.getByText("Software Engineer")).toBeInTheDocument();
		expect(screen.getByText("Acme Corp")).toBeInTheDocument();
		expect(screen.getByText(/Built awesome features/)).toBeInTheDocument();
	});

	it("renders work experience with Present end date", () => {
		const experience: WorkExperience = {
			id: "123",
			companyName: "Acme Corp",
			role: "Software Engineer",
			startDate: "2020-01-01",
			endDate: undefined,
			description: "Current role",
		};
		render(<WorkExperienceCard experience={experience} />);
		expect(screen.getByText(/Present/)).toBeInTheDocument();
	});

	it("renders work experience with credibility highlights", () => {
		const experience: WorkExperience = {
			id: "123",
			companyName: "Acme Corp",
			role: "Software Engineer",
			startDate: "2020-01-01",
			endDate: "2022-12-31",
			description: "Built features",
			credibilityHighlights: [
				{
					quote: "Excellent performer",
					sentiment: "POSITIVE" as const,
				},
			],
		};
		render(<WorkExperienceCard experience={experience} />);
		expect(screen.getByText(/Excellent performer/)).toBeInTheDocument();
	});

	it("handles missing description gracefully", () => {
		const experience: WorkExperience = {
			id: "123",
			companyName: "Acme Corp",
			role: "Software Engineer",
			startDate: "2020-01-01",
			endDate: "2022-12-31",
		};
		render(<WorkExperienceCard experience={experience} />);
		expect(screen.getByText("Software Engineer")).toBeInTheDocument();
		expect(screen.getByText("Acme Corp")).toBeInTheDocument();
	});
});

