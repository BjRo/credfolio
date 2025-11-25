import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { ReferenceLetterUpload } from "./ReferenceLetterUpload";
import * as referenceLettersApi from "@/lib/api/referenceLetters";

vi.mock("@/lib/api/referenceLetters", () => ({
	uploadReferenceLetter: vi.fn(),
}));

describe("ReferenceLetterUpload", () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	// T056: Unit test for ReferenceLetterUpload component when file selected shows file name
	it("when file selected, shows file name", async () => {
		// Arrange
		render(<ReferenceLetterUpload />);

		const file = new File(["test content"], "test-letter.txt", {
			type: "text/plain",
		});

		const input = document.querySelector(
			'input[type="file"]',
		) as HTMLInputElement;

		// Act
		fireEvent.change(input, { target: { files: [file] } });

		// Assert
		await waitFor(() => {
			expect(screen.getByText("test-letter.txt")).toBeInTheDocument();
		});
	});

	// T057: Unit test for ReferenceLetterUpload component when upload fails displays error message
	it("when upload fails, displays error message", async () => {
		// Arrange
		const mockUpload = vi.mocked(referenceLettersApi.uploadReferenceLetter);
		mockUpload.mockRejectedValue(new Error("Upload failed"));

		const onError = vi.fn();
		render(<ReferenceLetterUpload onError={onError} />);

		const file = new File(["test content"], "test-letter.txt", {
			type: "text/plain",
		});

		const input = document.querySelector(
			'input[type="file"]',
		) as HTMLInputElement;
		fireEvent.change(input, { target: { files: [file] } });

		// Act
		const uploadButton = await screen.findByRole("button", { name: /upload/i });
		fireEvent.click(uploadButton);

		// Assert
		await waitFor(() => {
			expect(onError).toHaveBeenCalledWith("Upload failed");
		});

		expect(screen.getByText("Upload failed")).toBeInTheDocument();
	});

	it("when invalid file type selected, shows error", async () => {
		// Arrange
		const onError = vi.fn();
		render(<ReferenceLetterUpload onError={onError} />);

		const file = new File(["pdf content"], "test-letter.pdf", {
			type: "application/pdf",
		});

		const input = document.querySelector(
			'input[type="file"]',
		) as HTMLInputElement;

		// Act
		fireEvent.change(input, { target: { files: [file] } });

		// Assert
		await waitFor(() => {
			expect(onError).toHaveBeenCalledWith("Please select a .txt or .md file");
		});
	});

	it("when upload succeeds, calls onUploadComplete callback", async () => {
		// Arrange
		const mockLetter = {
			id: "test-id",
			fileName: "test-letter.txt",
			uploadDate: new Date().toISOString(),
			status: "PENDING" as const,
		};

		const mockUpload = vi.mocked(referenceLettersApi.uploadReferenceLetter);
		mockUpload.mockResolvedValue(mockLetter);

		const onUploadComplete = vi.fn();
		render(<ReferenceLetterUpload onUploadComplete={onUploadComplete} />);

		const file = new File(["test content"], "test-letter.txt", {
			type: "text/plain",
		});

		const input = document.querySelector(
			'input[type="file"]',
		) as HTMLInputElement;
		fireEvent.change(input, { target: { files: [file] } });

		// Act
		const uploadButton = await screen.findByRole("button", { name: /upload/i });
		fireEvent.click(uploadButton);

		// Assert
		await waitFor(() => {
			expect(onUploadComplete).toHaveBeenCalledWith(mockLetter);
		});
	});
});
