import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import ExperienceList from "./ExperienceList";

describe("ExperienceList", () => {
	const mockCompanies = [
		{
			id: "c1",
			name: "Tech Corp",
			start_date: "2020-01-01",
			end_date: null, // Present
			roles: [
				{
					id: "r1",
					title: "Senior Dev",
					description: "Leading the team",
					is_verified: true,
					employer_feedback: "Great leader",
					skills: [
						{ id: "s1", name: "React" },
						{ id: "s2", name: "Go" },
					],
				},
			],
		},
		{
			id: "c2",
			name: "Old Job Inc",
			start_date: "2018-01-01",
			end_date: "2019-12-31",
			roles: [
				{
					id: "r2",
					title: "Junior Dev",
					description: "Fixing bugs",
					is_verified: false,
					skills: [{ id: "s3", name: "Java" }],
				},
			],
		},
	];

	it("renders company names and roles", () => {
		render(<ExperienceList companies={mockCompanies} />);
		expect(screen.getByText("Tech Corp")).toBeInTheDocument();
		expect(screen.getByText("Senior Dev")).toBeInTheDocument();
		expect(screen.getByText("Old Job Inc")).toBeInTheDocument();
		expect(screen.getByText("Junior Dev")).toBeInTheDocument();
	});

	it("renders dates correctly", () => {
		render(<ExperienceList companies={mockCompanies} />);
		// Use a function matcher to be resilient against whitespace/newlines in JSX
		expect(
			screen.getByText(
				(content) =>
					content.includes("Jan 2020") && content.includes("Present"),
			),
		).toBeInTheDocument();
		expect(
			screen.getByText(
				(content) =>
					content.includes("Jan 2018") && content.includes("Dec 2019"),
			),
		).toBeInTheDocument();
	});

	it("shows verified badge for verified roles", () => {
		render(<ExperienceList companies={mockCompanies} />);
		const badges = screen.getAllByText("Verified", { selector: "span" });
		expect(badges).toHaveLength(1); // Only for Senior Dev
	});

	it("renders employer feedback when present", () => {
		render(<ExperienceList companies={mockCompanies} />);
		expect(screen.getByText(/"Great leader"/)).toBeInTheDocument();
	});

	it("renders skills", () => {
		render(<ExperienceList companies={mockCompanies} />);
		expect(screen.getByText("React")).toBeInTheDocument();
		expect(screen.getByText("Go")).toBeInTheDocument();
		expect(screen.getByText("Java")).toBeInTheDocument();
	});

	it("renders empty state message", () => {
		render(<ExperienceList companies={[]} />);
		expect(screen.getByText("No experience recorded yet.")).toBeInTheDocument();
	});

	it("highlights skills when requested", () => {
		render(
			<ExperienceList companies={mockCompanies} highlightSkillIDs={["s1"]} />,
		);
		// React (s1) should have highlight classes (green bg)
		const reactSkill = screen.getByText("React");
		expect(reactSkill).toHaveClass("bg-green-200");

		// Go (s2) should have normal classes (gray bg)
		const goSkill = screen.getByText("Go");
		expect(goSkill).toHaveClass("bg-gray-200");
	});
});
