import { OpenAPI } from "./generated";

const API_BASE_URL =
	process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export interface ReferenceLetter {
	id: string;
	fileName: string;
	uploadDate: string;
	status: string;
}

/**
 * Upload a reference letter file
 */
export async function uploadReferenceLetter(
	file: File,
): Promise<ReferenceLetter> {
	const formData = new FormData();
	formData.append("file", file);

	const response = await fetch(`${API_BASE_URL}/reference-letters`, {
		method: "POST",
		body: formData,
		credentials: "include",
	});

	if (!response.ok) {
		const error = await response.text();
		throw new Error(error || "Failed to upload reference letter");
	}

	return response.json();
}

/**
 * Get all reference letters for the current user
 */
export async function getReferenceLetters(): Promise<ReferenceLetter[]> {
	const response = await fetch(`${API_BASE_URL}/reference-letters`, {
		method: "GET",
		credentials: "include",
	});

	if (!response.ok) {
		throw new Error("Failed to get reference letters");
	}

	return response.json();
}
