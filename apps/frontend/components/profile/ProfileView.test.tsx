import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { ProfileView } from "./ProfileView";
import * as profileApi from "@/lib/api/profile";

vi.mock("@/lib/api/profile", () => ({
	getProfile: vi.fn(),
}));

describe("ProfileView", () => {
	const mockProfile = {
		id: "test-id",
		summary:
			"Experienced software engineer with expertise in Go and TypeScript",
		workExperiences: [
			{
				id: "exp-1",
				companyName: "Tech Corp",
				role: "Senior Engineer",
				startDate: "2020-01-01",
				endDate: "2023-12-31",
				description: "Led development of microservices",
				credibilityHighlights: [
					{
						quote: "Outstanding technical leadership",
						sentiment: "POSITIVE" as const,
					},
				],
			},
			{
				id: "exp-2",
				companyName: "Startup Inc",
				role: "Software Engineer",
				startDate: "2018-06-01",
				endDate: "2019-12-31",
				description: "Full-stack development",
				credibilityHighlights: [],
			},
		],
		skills: ["Go", "TypeScript", "React", "PostgreSQL"],
	};

	beforeEach(() => {
		vi.clearAllMocks();
	});

	// T088: Unit test for ProfileView component when profile loaded displays all sections
	it("when profile loaded, displays all sections", async () => {
		// Arrange
		const mockGetProfile = vi.mocked(profileApi.getProfile);
		mockGetProfile.mockResolvedValue(mockProfile);

		// Act
		render(<ProfileView />);

		// Assert
		await waitFor(() => {
			// Summary section
			expect(
				screen.getByText(/experienced software engineer/i),
			).toBeInTheDocument();

			// Work experiences
			expect(screen.getByText("Tech Corp")).toBeInTheDocument();
			expect(screen.getByText("Startup Inc")).toBeInTheDocument();

			// Skills
			expect(screen.getByText("Go")).toBeInTheDocument();
			expect(screen.getByText("TypeScript")).toBeInTheDocument();

			// Credibility highlight
			expect(
				screen.getByText(/outstanding technical leadership/i),
			).toBeInTheDocument();
		});
	});

	// T089: Unit test for ProfileView component when loading shows loading state
	it("when loading, shows loading state", () => {
		// Arrange
		const mockGetProfile = vi.mocked(profileApi.getProfile);
		mockGetProfile.mockImplementation(
			() => new Promise(() => {}), // Never resolves
		);

		// Act
		render(<ProfileView />);

		// Assert
		expect(screen.getByText(/loading/i)).toBeInTheDocument();
	});

	it("when profile has no summary, shows placeholder", async () => {
		// Arrange
		const profileNoSummary = { ...mockProfile, summary: "" };
		const mockGetProfile = vi.mocked(profileApi.getProfile);
		mockGetProfile.mockResolvedValue(profileNoSummary);

		// Act
		render(<ProfileView />);

		// Assert
		await waitFor(() => {
			expect(screen.getByText(/no summary available/i)).toBeInTheDocument();
		});
	});

	it("when profile has no work experiences, shows placeholder", async () => {
		// Arrange
		const profileNoExp = { ...mockProfile, workExperiences: [] };
		const mockGetProfile = vi.mocked(profileApi.getProfile);
		mockGetProfile.mockResolvedValue(profileNoExp);

		// Act
		render(<ProfileView />);

		// Assert
		await waitFor(() => {
			expect(screen.getByText(/no work experience yet/i)).toBeInTheDocument();
		});
	});

	it("when fetch fails, shows error message", async () => {
		// Arrange
		const mockGetProfile = vi.mocked(profileApi.getProfile);
		mockGetProfile.mockRejectedValue(new Error("Failed to fetch profile"));

		// Act
		render(<ProfileView />);

		// Assert
		await waitFor(() => {
			expect(screen.getByText(/failed to load profile/i)).toBeInTheDocument();
		});
	});
});
