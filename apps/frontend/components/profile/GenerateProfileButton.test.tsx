import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { GenerateProfileButton } from "./GenerateProfileButton";
import * as profileApi from "@/lib/api/profile";

vi.mock("@/lib/api/profile", () => ({
	generateProfile: vi.fn(),
}));

describe("GenerateProfileButton", () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	// T058: Unit test for GenerateProfileButton component when clicked triggers generation
	it("when clicked, triggers generation", async () => {
		// Arrange
		const mockProfile = {
			id: "test-profile-id",
			summary: "Test summary",
			workExperiences: [],
			skills: [],
		};

		const mockGenerate = vi.mocked(profileApi.generateProfile);
		mockGenerate.mockResolvedValue(mockProfile);

		const onGenerateComplete = vi.fn();
		render(<GenerateProfileButton onGenerateComplete={onGenerateComplete} />);

		// Act
		const button = screen.getByRole("button", {
			name: /generate profile with ai/i,
		});
		fireEvent.click(button);

		// Assert
		await waitFor(() => {
			expect(mockGenerate).toHaveBeenCalled();
		});

		await waitFor(() => {
			expect(onGenerateComplete).toHaveBeenCalledWith(mockProfile);
		});
	});

	it("when generating, shows loading state", async () => {
		// Arrange
		const mockGenerate = vi.mocked(profileApi.generateProfile);
		mockGenerate.mockImplementation(
			() =>
				new Promise((resolve) =>
					setTimeout(
						() =>
							resolve({
								id: "test-id",
								summary: "",
								workExperiences: [],
								skills: [],
							}),
						100,
					),
				),
		);

		render(<GenerateProfileButton />);

		// Act
		const button = screen.getByRole("button", {
			name: /generate profile with ai/i,
		});
		fireEvent.click(button);

		// Assert
		await waitFor(() => {
			expect(screen.getByText(/generating profile/i)).toBeInTheDocument();
		});
	});

	it("when generation fails, calls onError callback", async () => {
		// Arrange
		const mockGenerate = vi.mocked(profileApi.generateProfile);
		mockGenerate.mockRejectedValue(new Error("Generation failed"));

		const onError = vi.fn();
		render(<GenerateProfileButton onError={onError} />);

		// Act
		const button = screen.getByRole("button", {
			name: /generate profile with ai/i,
		});
		fireEvent.click(button);

		// Assert
		await waitFor(() => {
			expect(onError).toHaveBeenCalledWith("Generation failed");
		});
	});

	it("when disabled, button is not clickable", () => {
		// Arrange
		render(<GenerateProfileButton disabled />);

		// Act & Assert
		const button = screen.getByRole("button", {
			name: /generate profile with ai/i,
		});
		expect(button).toBeDisabled();
	});
});
