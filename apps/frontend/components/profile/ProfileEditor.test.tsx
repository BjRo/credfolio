import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import ProfileEditor from "./ProfileEditor";
import { DefaultService } from "../../lib/api/generated";

vi.mock("../../lib/api/generated", () => ({
	DefaultService: {
		updateProfile: vi.fn(),
	},
}));

describe("ProfileEditor", () => {
	const mockProfile = {
		summary: "Original Summary",
	};

	it("renders with initial profile data", () => {
		render(<ProfileEditor profile={mockProfile} onUpdate={() => {}} />);
		expect(screen.getByDisplayValue("Original Summary")).toBeInTheDocument();
	});

	it("updates profile summary", async () => {
		const onUpdate = vi.fn();
		const updatedProfile = { ...mockProfile, summary: "Updated Summary" };
		vi.mocked(DefaultService.updateProfile).mockResolvedValueOnce(
			updatedProfile as never,
		);

		// Mock window.alert since it's used in the component
		const alertMock = vi.spyOn(window, "alert").mockImplementation(() => {});

		render(<ProfileEditor profile={mockProfile} onUpdate={onUpdate} />);

		const textarea = screen.getByDisplayValue("Original Summary");
		fireEvent.change(textarea, { target: { value: "Updated Summary" } });

		const saveButton = screen.getByRole("button", { name: /Save Changes/i });
		fireEvent.click(saveButton);

		expect(saveButton).toBeDisabled();
		expect(saveButton).toHaveTextContent("Saving...");

		await waitFor(() => {
			expect(DefaultService.updateProfile).toHaveBeenCalledWith({
				summary: "Updated Summary",
			});
			expect(onUpdate).toHaveBeenCalledWith(updatedProfile);
			expect(alertMock).toHaveBeenCalledWith("Profile updated!");
		});

		alertMock.mockRestore();
	});

	it("handles update failure", async () => {
		vi.mocked(DefaultService.updateProfile).mockRejectedValueOnce(
			new Error("Update failed"),
		);
		const alertMock = vi.spyOn(window, "alert").mockImplementation(() => {});
		const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

		render(<ProfileEditor profile={mockProfile} onUpdate={() => {}} />);

		const saveButton = screen.getByRole("button", { name: /Save Changes/i });
		fireEvent.click(saveButton);

		await waitFor(() => {
			expect(alertMock).toHaveBeenCalledWith("Failed to update profile.");
		});

		alertMock.mockRestore();
		consoleSpy.mockRestore();
	});
});
