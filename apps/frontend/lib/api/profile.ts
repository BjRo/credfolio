import { DefaultService } from "./generated";
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
