import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import DownloadCVButton from "./DownloadCVButton";
import * as profileApi from "../../lib/api/profile";

vi.mock("../../lib/api/profile", () => ({
	downloadCV: vi.fn(),
}));

describe("DownloadCVButton", () => {
	const mockCreateObjectURL = vi.fn(() => "blob:mock-url");
	const mockRevokeObjectURL = vi.fn();
	const mockClick = vi.fn();

	beforeEach(() => {
		// Mock window.URL methods
		global.URL.createObjectURL = mockCreateObjectURL;
		global.URL.revokeObjectURL = mockRevokeObjectURL;

		// Mock anchor element click - spy on the actual element after it's created
		vi.spyOn(HTMLAnchorElement.prototype, "click").mockImplementation(
			mockClick,
		);

		// Suppress console errors in tests
		vi.spyOn(console, "error").mockImplementation(() => {});
	});

	afterEach(() => {
		vi.clearAllMocks();
		vi.restoreAllMocks();
	});

	it("renders the button correctly", () => {
		render(<DownloadCVButton profileId="test-profile-id" />);
		expect(
			screen.getByRole("button", { name: /Download CV/i }),
		).toBeInTheDocument();
	});

	it("renders with custom className", () => {
		const { container } = render(
			<DownloadCVButton profileId="test-id" className="custom-class" />,
		);
		const button = container.querySelector("button");
		expect(button?.className).toContain("custom-class");
	});

	it("triggers CV download on click without jobMatchId", async () => {
		const mockBlob = new Blob(["test"], { type: "application/pdf" });
		vi.mocked(profileApi.downloadCV).mockResolvedValueOnce(mockBlob);

		render(<DownloadCVButton profileId="test-profile-id" />);

		const button = screen.getByRole("button", { name: /Download CV/i });
		fireEvent.click(button);

		expect(button).toBeDisabled();
		expect(screen.getByText(/Downloading.../i)).toBeInTheDocument();

		await waitFor(() => {
			expect(profileApi.downloadCV).toHaveBeenCalledWith(
				"test-profile-id",
				undefined,
			);
		});

		// Verify blob URL creation and cleanup
		expect(mockCreateObjectURL).toHaveBeenCalledWith(mockBlob);
		expect(mockClick).toHaveBeenCalled();
		// Give a moment for cleanup
		await waitFor(
			() => {
				expect(mockRevokeObjectURL).toHaveBeenCalledWith("blob:mock-url");
			},
			{ timeout: 100 },
		);

		expect(button).not.toBeDisabled();
		expect(screen.getByText(/Download CV/i)).toBeInTheDocument();
	});

	it("triggers tailored CV download on click with jobMatchId", async () => {
		const mockBlob = new Blob(["test"], { type: "application/pdf" });
		vi.mocked(profileApi.downloadCV).mockResolvedValueOnce(mockBlob);

		render(
			<DownloadCVButton
				profileId="test-profile-id"
				jobMatchId="test-job-match-id"
			/>,
		);

		const button = screen.getByRole("button", { name: /Download CV/i });
		fireEvent.click(button);

		await waitFor(() => {
			expect(profileApi.downloadCV).toHaveBeenCalledWith(
				"test-profile-id",
				"test-job-match-id",
			);
		});

		// Verify download was triggered
		expect(mockCreateObjectURL).toHaveBeenCalled();
		expect(mockClick).toHaveBeenCalled();
	});

	it("handles download error", async () => {
		vi.mocked(profileApi.downloadCV).mockRejectedValueOnce(
			new Error("Download failed"),
		);

		render(<DownloadCVButton profileId="test-profile-id" />);

		const button = screen.getByRole("button", { name: /Download CV/i });
		fireEvent.click(button);

		await waitFor(
			() => {
				expect(screen.getByText(/Download failed/i)).toBeInTheDocument();
			},
			{ timeout: 2000 },
		);

		expect(button).not.toBeDisabled();
	});

	it("handles download error with error message", async () => {
		vi.mocked(profileApi.downloadCV).mockRejectedValueOnce(
			new Error("Network error"),
		);

		render(<DownloadCVButton profileId="test-profile-id" />);

		const button = screen.getByRole("button", { name: /Download CV/i });
		fireEvent.click(button);

		await waitFor(() => {
			expect(screen.getByText(/Network error/i)).toBeInTheDocument();
		});
	});

	it("handles download error with non-Error object", async () => {
		vi.mocked(profileApi.downloadCV).mockRejectedValueOnce("Unknown error");

		render(<DownloadCVButton profileId="test-profile-id" />);

		const button = screen.getByRole("button", { name: /Download CV/i });
		fireEvent.click(button);

		await waitFor(() => {
			expect(
				screen.getByText(/An unexpected error occurred/i),
			).toBeInTheDocument();
		});
	});
});
