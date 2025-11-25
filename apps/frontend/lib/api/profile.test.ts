import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { generateProfile, getProfile, updateProfile } from "./profile";

describe("profile API", () => {
	const originalFetch = global.fetch;

	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		global.fetch = originalFetch;
	});

	// T061: Unit test for API client when generating profile calls correct endpoint
	describe("generateProfile", () => {
		it("when generating profile, calls correct endpoint", async () => {
			// Arrange
			const mockProfile = {
				id: "test-id",
				summary: "Test summary",
				workExperiences: [],
				skills: [],
			};

			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: () => Promise.resolve(mockProfile),
			});
			global.fetch = mockFetch;

			// Act
			await generateProfile();

			// Assert
			expect(mockFetch).toHaveBeenCalled();
			const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
			expect(url).toContain("/profile/generate");
			expect(options.method).toBe("POST");
		});

		it("returns profile data on success", async () => {
			// Arrange
			const mockProfile = {
				id: "test-id",
				summary: "Test summary",
				workExperiences: [
					{
						id: "exp-1",
						companyName: "Acme Corp",
						role: "Engineer",
						startDate: "2020-01-01",
						description: "Work",
						credibilityHighlights: [],
					},
				],
				skills: ["TypeScript"],
			};

			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: () => Promise.resolve(mockProfile),
			});
			global.fetch = mockFetch;

			// Act
			const result = await generateProfile();

			// Assert
			expect(result).toEqual(mockProfile);
		});

		it("when generation fails, throws error", async () => {
			// Arrange
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 500,
				statusText: "Internal Server Error",
			});
			global.fetch = mockFetch;

			// Act & Assert
			await expect(generateProfile()).rejects.toThrow();
		});
	});

	describe("getProfile", () => {
		it("when getting profile, calls correct endpoint", async () => {
			// Arrange
			const mockProfile = {
				id: "test-id",
				summary: "Test summary",
				workExperiences: [],
				skills: [],
			};

			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: () => Promise.resolve(mockProfile),
			});
			global.fetch = mockFetch;

			// Act
			await getProfile();

			// Assert
			expect(mockFetch).toHaveBeenCalled();
			const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
			expect(url).toContain("/profile");
			expect(options.method).toBe("GET");
		});

		it("when profile not found, throws error", async () => {
			// Arrange
			const mockFetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 404,
				statusText: "Not Found",
			});
			global.fetch = mockFetch;

			// Act & Assert
			await expect(getProfile()).rejects.toThrow("Profile not found");
		});
	});

	describe("updateProfile", () => {
		it("when updating profile, sends JSON data", async () => {
			// Arrange
			const mockProfile = {
				id: "test-id",
				summary: "Updated summary",
				workExperiences: [],
				skills: [],
			};

			const mockFetch = vi.fn().mockResolvedValue({
				ok: true,
				json: () => Promise.resolve(mockProfile),
			});
			global.fetch = mockFetch;

			// Act
			await updateProfile({ summary: "Updated summary" });

			// Assert
			expect(mockFetch).toHaveBeenCalled();
			const [url, options] = mockFetch.mock.calls[0] as [string, RequestInit];
			expect(url).toContain("/profile");
			expect(options.method).toBe("PUT");
			expect(options.headers).toEqual(
				expect.objectContaining({
					"Content-Type": "application/json",
				}),
			);
			expect(JSON.parse(options.body as string)).toEqual({
				summary: "Updated summary",
			});
		});
	});
});
