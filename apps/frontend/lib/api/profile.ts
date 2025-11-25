const API_BASE_URL =
	process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export interface CredibilityHighlight {
	quote: string;
	sentiment: "POSITIVE" | "NEUTRAL";
}

export interface WorkExperience {
	id: string;
	companyName: string;
	role: string;
	startDate: string;
	endDate?: string;
	description: string;
	credibilityHighlights?: CredibilityHighlight[];
}

export interface Profile {
	id: string;
	summary: string;
	workExperiences: WorkExperience[];
	skills: string[];
}

export interface JobMatch {
	id: string;
	matchScore: number;
	tailoredSummary: string;
}

/**
 * Get the current user's profile
 */
export async function getProfile(): Promise<Profile> {
	const response = await fetch(`${API_BASE_URL}/profile`, {
		method: "GET",
		credentials: "include",
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error("Profile not found");
		}
		throw new Error("Failed to get profile");
	}

	return response.json();
}

/**
 * Generate profile from uploaded reference letters
 */
export async function generateProfile(): Promise<Profile> {
	const response = await fetch(`${API_BASE_URL}/profile/generate`, {
		method: "POST",
		credentials: "include",
	});

	if (!response.ok) {
		const error = await response.text();
		throw new Error(error || "Failed to generate profile");
	}

	return response.json();
}

/**
 * Update the current user's profile
 */
export async function updateProfile(data: {
	summary: string;
}): Promise<Profile> {
	const response = await fetch(`${API_BASE_URL}/profile`, {
		method: "PUT",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(data),
		credentials: "include",
	});

	if (!response.ok) {
		throw new Error("Failed to update profile");
	}

	return response.json();
}

/**
 * Tailor profile to a job description
 */
export async function tailorProfile(jobDescription: string): Promise<JobMatch> {
	const response = await fetch(`${API_BASE_URL}/profile/tailor`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ jobDescription }),
		credentials: "include",
	});

	if (!response.ok) {
		throw new Error("Failed to tailor profile");
	}

	return response.json();
}

/**
 * Download CV as PDF
 */
export async function downloadCV(profileId: string): Promise<Blob> {
	const response = await fetch(`${API_BASE_URL}/profile/${profileId}/cv`, {
		method: "GET",
		credentials: "include",
	});

	if (!response.ok) {
		throw new Error("Failed to download CV");
	}

	return response.blob();
}
