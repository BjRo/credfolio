import type { UserProfile, TailoringResult } from "../types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function uploadReferenceLetter(
	file: File,
): Promise<{ status: string; user_id: string }> {
	const formData = new FormData();
	formData.append("file", file);

	const res = await fetch(`${API_URL}/api/upload`, {
		method: "POST",
		body: formData,
	});

	if (!res.ok) {
		const errorText = await res.text();
		throw new Error(errorText || "Upload failed");
	}

	return res.json();
}

export async function getProfile(userID: string): Promise<UserProfile> {
	const res = await fetch(`${API_URL}/api/profile?user_id=${userID}`);
	if (!res.ok) {
		throw new Error("Failed to fetch profile");
	}
	return res.json();
}

export async function tailorProfile(
	userID: string,
	jobDescription: string,
): Promise<TailoringResult> {
	const res = await fetch(`${API_URL}/api/profile/tailor`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({
			user_id: userID,
			job_description: jobDescription,
		}),
	});

	if (!res.ok) {
		const errorText = await res.text();
		throw new Error(errorText || "Tailoring failed");
	}

	return res.json();
}
