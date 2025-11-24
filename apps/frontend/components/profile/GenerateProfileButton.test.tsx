import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import GenerateProfileButton from "./GenerateProfileButton";
import * as profileApi from "../../lib/api/profile";

vi.mock("../../lib/api/profile", () => ({
	generateProfile: vi.fn(),
}));

describe("GenerateProfileButton", () => {
	it("renders the button correctly", () => {
		render(<GenerateProfileButton />);
		expect(
			screen.getByRole("button", { name: /Generate Smart Profile/i }),
		).toBeInTheDocument();
	});

	it("triggers profile generation on click", async () => {
		const onGenerateComplete = vi.fn();
		const mockProfile = { id: "123", summary: "Test Profile" };
		vi.mocked(profileApi.generateProfile).mockResolvedValueOnce(
			mockProfile as never,
		);

		render(<GenerateProfileButton onGenerateComplete={onGenerateComplete} />);

		const button = screen.getByRole("button", {
			name: /Generate Smart Profile/i,
		});
		fireEvent.click(button);

		expect(button).toBeDisabled();
		expect(button).toHaveTextContent("Generating Profile...");

		await waitFor(() => {
			expect(profileApi.generateProfile).toHaveBeenCalled();
			expect(onGenerateComplete).toHaveBeenCalledWith(mockProfile);
			expect(button).not.toBeDisabled();
			expect(button).toHaveTextContent("Generate Smart Profile");
		});
	});

	it("handles generation error", async () => {
		const consoleErrorSpy = vi
			.spyOn(console, "error")
			.mockImplementation(() => {});
		vi.mocked(profileApi.generateProfile).mockRejectedValueOnce(
			new Error("Generation failed"),
		);

		render(<GenerateProfileButton />);

		const button = screen.getByRole("button", {
			name: /Generate Smart Profile/i,
		});
		fireEvent.click(button);

		await waitFor(() => {
			expect(screen.getByText(/Generation failed/i)).toBeInTheDocument();
			expect(button).not.toBeDisabled();
		});

		consoleErrorSpy.mockRestore();
	});
});
