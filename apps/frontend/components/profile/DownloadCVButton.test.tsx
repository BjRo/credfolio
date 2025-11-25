import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { DownloadCVButton } from "./DownloadCVButton";
import * as profileApi from "@/lib/api/profile";

vi.mock("@/lib/api/profile", () => ({
	downloadCV: vi.fn(),
}));

describe("DownloadCVButton", () => {
	beforeEach(() => {
		vi.clearAllMocks();
		// Mock URL.createObjectURL and URL.revokeObjectURL
		global.URL.createObjectURL = vi.fn(() => "blob:mock-url");
		global.URL.revokeObjectURL = vi.fn();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	// T134: Unit test for DownloadCVButton component when clicked triggers download
	it("when clicked, triggers download", async () => {
		// Arrange
		const mockBlob = new Blob(["mock pdf content"], {
			type: "application/pdf",
		});
		const mockDownloadCV = vi.mocked(profileApi.downloadCV);
		mockDownloadCV.mockResolvedValue(mockBlob);

		// Act
		render(<DownloadCVButton profileId="test-profile-id" />);
		const button = screen.getByRole("button", { name: /download/i });
		fireEvent.click(button);

		// Assert
		await waitFor(() => {
			expect(mockDownloadCV).toHaveBeenCalledWith("test-profile-id", undefined);
		});
	});

	// T135: Unit test for DownloadCVButton component when download fails displays error
	it("when download fails, displays error", async () => {
		// Arrange
		const mockDownloadCV = vi.mocked(profileApi.downloadCV);
		mockDownloadCV.mockRejectedValue(new Error("Download failed"));

		const onError = vi.fn();

		// Act
		render(<DownloadCVButton profileId="test-profile-id" onError={onError} />);
		const button = screen.getByRole("button", { name: /download/i });
		fireEvent.click(button);

		// Assert
		await waitFor(() => {
			expect(onError).toHaveBeenCalledWith("Failed to download CV");
		});
	});

	it("shows loading state while downloading", async () => {
		// Arrange
		const mockDownloadCV = vi.mocked(profileApi.downloadCV);
		mockDownloadCV.mockImplementation(
			() =>
				new Promise((resolve) => setTimeout(() => resolve(new Blob()), 100)),
		);

		// Act
		render(<DownloadCVButton profileId="test-profile-id" />);
		const button = screen.getByRole("button", { name: /download/i });
		fireEvent.click(button);

		// Assert
		await waitFor(() => {
			expect(screen.getByText(/downloading/i)).toBeInTheDocument();
		});
	});

	it("renders with default text", () => {
		// Arrange & Act
		render(<DownloadCVButton profileId="test-profile-id" />);

		// Assert
		expect(
			screen.getByRole("button", { name: /download cv/i }),
		).toBeInTheDocument();
	});

	it("can be disabled", () => {
		// Arrange & Act
		render(<DownloadCVButton profileId="test-profile-id" disabled />);

		// Assert
		expect(screen.getByRole("button")).toBeDisabled();
	});
});
