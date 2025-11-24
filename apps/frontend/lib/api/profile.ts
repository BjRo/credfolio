import { DefaultService } from "./generated";

export const generateProfile = async () => {
	return DefaultService.generateProfile();
};

export const getProfile = async () => {
	return DefaultService.getProfile();
};
