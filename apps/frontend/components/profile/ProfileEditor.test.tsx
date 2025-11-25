import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { ProfileEditor } from "./ProfileEditor";
import * as profileApi from "@/lib/api/profile";

vi.mock("@/lib/api/profile", () => ({
	updateProfile: vi.fn(),
}));

describe("ProfileEditor", () => {
	const mockProfile = {
		id: "test-profile-id",
		summary: "Original summary",
		workExperiences: [
			{
				id: "exp-1",
				companyName: "Acme Corp",
				role: "Software Engineer",
				startDate: "2020-01-01",
				endDate: "2023-12-31",
				description: "Built software",
				credibilityHighlights: [
					{
						quote: "Exceptional team player",
						sentiment: "POSITIVE" as const,
					},
				],
			},
		],
		skills: ["TypeScript", "React", "Node.js"],
	};

	beforeEach(() => {
		vi.clearAllMocks();
	});

	// T059: Unit test for ProfileEditor component when editing field updates value
	it("when editing field, updates value", async () => {
		// Arrange
		render(<ProfileEditor profile={mockProfile} />);

		// Act - click edit button for summary
		const editButton = screen.getByRole("button", { name: /edit/i });
		fireEvent.click(editButton);

		// Find the textarea and update it
		const textarea = screen.getByPlaceholderText(
			/write a professional summary/i,
		);
		fireEvent.change(textarea, { target: { value: "Updated summary" } });

		// Assert
		expect(textarea).toHaveValue("Updated summary");
	});

	it("displays profile summary", () => {
		// Arrange & Act
		render(<ProfileEditor profile={mockProfile} />);

		// Assert
		expect(screen.getByText("Original summary")).toBeInTheDocument();
	});

	it("displays work experiences", () => {
		// Arrange & Act
		render(<ProfileEditor profile={mockProfile} />);

		// Assert
		expect(screen.getByText("Acme Corp")).toBeInTheDocument();
		expect(screen.getByText("Software Engineer")).toBeInTheDocument();
	});

	it("displays skills", () => {
		// Arrange & Act
		render(<ProfileEditor profile={mockProfile} />);

		// Assert
		expect(screen.getByText("TypeScript")).toBeInTheDocument();
		expect(screen.getByText("React")).toBeInTheDocument();
		expect(screen.getByText("Node.js")).toBeInTheDocument();
	});

	it("when save clicked, calls updateProfile API", async () => {
		// Arrange
		const mockUpdate = vi.mocked(profileApi.updateProfile);
		mockUpdate.mockResolvedValue({
			...mockProfile,
			summary: "Updated summary",
		});

		const onUpdate = vi.fn();
		render(<ProfileEditor profile={mockProfile} onUpdate={onUpdate} />);

		// Act - click edit, update, and save
		const editButton = screen.getByRole("button", { name: /edit/i });
		fireEvent.click(editButton);

		const textarea = screen.getByPlaceholderText(
			/write a professional summary/i,
		);
		fireEvent.change(textarea, { target: { value: "Updated summary" } });

		const saveButton = screen.getByRole("button", { name: /save/i });
		fireEvent.click(saveButton);

		// Assert
		await waitFor(() => {
			expect(mockUpdate).toHaveBeenCalledWith({ summary: "Updated summary" });
		});
	});

	it("when cancel clicked, reverts changes", async () => {
		// Arrange
		render(<ProfileEditor profile={mockProfile} />);

		// Act - click edit, update, then cancel
		const editButton = screen.getByRole("button", { name: /edit/i });
		fireEvent.click(editButton);

		const textarea = screen.getByPlaceholderText(
			/write a professional summary/i,
		);
		fireEvent.change(textarea, { target: { value: "Changed summary" } });

		const cancelButton = screen.getByRole("button", { name: /cancel/i });
		fireEvent.click(cancelButton);

		// Assert - should show original summary again
		expect(screen.getByText("Original summary")).toBeInTheDocument();
	});

	it("displays credibility highlights when work experience expanded", async () => {
		// Arrange
		render(<ProfileEditor profile={mockProfile} />);

		// Act - find and click the expand button
		const expandButton = screen.getByRole("button", {
			name: /expand details|collapse details/i,
		});
		fireEvent.click(expandButton);

		// Assert
		await waitFor(() => {
			expect(screen.getByText(/exceptional team player/i)).toBeInTheDocument();
		});
	});

	it("displays message when no summary exists", () => {
		// Arrange
		const profileNoSummary = { ...mockProfile, summary: "" };

		// Act
		render(<ProfileEditor profile={profileNoSummary} />);

		// Assert
		expect(
			screen.getByText(/no summary yet. click edit to add one/i),
		).toBeInTheDocument();
	});
});
