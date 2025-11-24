import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import ReferenceLetterUpload from "./ReferenceLetterUpload";
import * as referenceLettersApi from "../../lib/api/referenceLetters";

// Mock the API module
vi.mock("../../lib/api/referenceLetters", () => ({
	uploadReferenceLetter: vi.fn(),
}));

describe("ReferenceLetterUpload", () => {
	it("renders the upload component correctly", () => {
		render(<ReferenceLetterUpload />);
		expect(screen.getByText(/Upload Reference Letter/i)).toBeInTheDocument();
		expect(screen.getByText(/Choose file/i)).toBeInTheDocument();
	});

	it("handles file selection and upload success", async () => {
		const onUploadComplete = vi.fn();
		vi.mocked(referenceLettersApi.uploadReferenceLetter).mockResolvedValueOnce(
			{} as never,
		);

		render(<ReferenceLetterUpload onUploadComplete={onUploadComplete} />);

		const file = new File(["dummy content"], "test.pdf", {
			type: "application/pdf",
		});
		const input = screen.getByLabelText(/Choose file/i);

		fireEvent.change(input, { target: { files: [file] } });

		expect(screen.getByText(/Uploading.../i)).toBeInTheDocument();

		await waitFor(() => {
			expect(referenceLettersApi.uploadReferenceLetter).toHaveBeenCalledWith(
				file,
			);
			expect(
				screen.getByText(/Reference letter uploaded successfully!/i),
			).toBeInTheDocument();
			expect(onUploadComplete).toHaveBeenCalled();
		});
	});

	it("handles upload failure", async () => {
		const consoleErrorSpy = vi
			.spyOn(console, "error")
			.mockImplementation(() => {});
		vi.mocked(referenceLettersApi.uploadReferenceLetter).mockRejectedValueOnce(
			new Error("Upload failed"),
		);

		render(<ReferenceLetterUpload />);

		const file = new File(["dummy content"], "test.pdf", {
			type: "application/pdf",
		});
		const input = screen.getByLabelText(/Choose file/i);

		fireEvent.change(input, { target: { files: [file] } });

		await waitFor(() => {
			expect(
				screen.getByText(/Failed to upload reference letter/i),
			).toBeInTheDocument();
		});

		consoleErrorSpy.mockRestore();
	});
});
