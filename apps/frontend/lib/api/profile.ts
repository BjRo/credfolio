import { DefaultService, OpenAPI } from "./generated";
import type { JobMatch } from "./generated/models/JobMatch";

export const generateProfile = async () => {
	return DefaultService.generateProfile();
};

export const getProfile = async () => {
	return DefaultService.getProfile();
};

export const tailorProfile = async (
	jobDescription: string,
): Promise<JobMatch> => {
	return DefaultService.tailorProfile({
		jobDescription,
	});
};

export const downloadCV = async (
	profileId: string,
	jobMatchId?: string,
): Promise<Blob> => {
	// Construct URL with query parameter if needed
	let url = `${OpenAPI.BASE || ""}/profile/${profileId}/cv`;
	if (jobMatchId) {
		url += `?jobMatchId=${encodeURIComponent(jobMatchId)}`;
	}

	const response = await fetch(url, {
		method: "GET",
		credentials: OpenAPI.CREDENTIALS,
		headers: (OpenAPI.HEADERS as HeadersInit) || {},
	});

	if (!response.ok) {
		throw new Error(`Failed to download CV: ${response.statusText}`);
	}

	return response.blob();
};
