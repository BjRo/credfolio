import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { uploadReferenceLetter } from "./referenceLetters";

describe("referenceLetters API", () => {
	const originalFetch = global.fetch;

	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		global.fetch = originalFetch;
	});

	// T060: Unit test for API client when posting reference letter sends multipart form data
	describe("uploadReferenceLetter", () => {
		it("when posting reference letter, sends multipart form data", async () => {
			// Arrange
			const mockResponse = {
				id: "test-id",
				fileName: "test.txt",
				uploadDate: new Date().toISOString(),
				status: "PENDING",
			};

			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: () => Promise.resolve(mockResponse),
			});
			global.fetch = mockFetch;

			const file = new File(["test content"], "test.txt", {
				type: "text/plain",
			});

			// Act
			await uploadReferenceLetter(file);

			// Assert
			expect(mockFetch).toHaveBeenCalled();
			const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
			expect(url).toContain("/reference-letters");
			expect(options.method).toBe("POST");
			expect(options.body).toBeInstanceOf(FormData);

			const formData = options.body as FormData;
			expect(formData.get("file")).toBeInstanceOf(File);
		});

		it("when upload fails, throws error", async () => {
			// Arrange
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 400,
				statusText: "Bad Request",
			});
			global.fetch = mockFetch;

			const file = new File(["test content"], "test.txt", {
				type: "text/plain",
			});

			// Act & Assert
			await expect(uploadReferenceLetter(file)).rejects.toThrow();
		});

		it("returns reference letter data on success", async () => {
			// Arrange
			const mockResponse = {
				id: "test-id",
				fileName: "test.txt",
				uploadDate: "2024-01-01T00:00:00Z",
				status: "PENDING" as const,
			};

			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: () => Promise.resolve(mockResponse),
			});
			global.fetch = mockFetch;

			const file = new File(["test content"], "test.txt", {
				type: "text/plain",
			});

			// Act
			const result = await uploadReferenceLetter(file);

			// Assert
			expect(result).toEqual(mockResponse);
		});
	});
});
