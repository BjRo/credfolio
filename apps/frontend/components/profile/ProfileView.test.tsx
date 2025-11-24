import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import ProfileView from "./ProfileView";
import type { Profile } from "../../lib/api/generated/models/Profile";

describe("ProfileView", () => {
	it("renders profile with summary", () => {
		const profile: Profile = {
			id: "123",
			summary: "Experienced software engineer",
		};
		render(<ProfileView profile={profile} />);
		expect(screen.getByText("Professional Profile")).toBeInTheDocument();
		expect(
			screen.getByText("Experienced software engineer"),
		).toBeInTheDocument();
	});

	it("renders profile with skills", () => {
		const profile: Profile = {
			id: "123",
			summary: "Test summary",
			skills: ["JavaScript", "TypeScript"],
		};
		render(<ProfileView profile={profile} />);
		expect(screen.getByText("Skills")).toBeInTheDocument();
		expect(screen.getByText("JavaScript")).toBeInTheDocument();
		expect(screen.getByText("TypeScript")).toBeInTheDocument();
	});

	it("renders profile with work experiences", () => {
		const profile: Profile = {
			id: "123",
			summary: "Test summary",
			workExperiences: [
				{
					id: "we1",
					companyName: "Acme Corp",
					role: "Software Engineer",
					startDate: "2020-01-01",
					endDate: "2022-12-31",
					description: "Built features",
				},
			],
		};
		render(<ProfileView profile={profile} />);
		expect(screen.getByText("Work Experience")).toBeInTheDocument();
		expect(screen.getByText("Software Engineer")).toBeInTheDocument();
		expect(screen.getByText("Acme Corp")).toBeInTheDocument();
	});

	it("renders empty state when no data", () => {
		const profile: Profile = {
			id: "123",
		};
		render(<ProfileView profile={profile} />);
		expect(
			screen.getByText(/No profile data available yet/),
		).toBeInTheDocument();
	});
});
