/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { JobMatch } from "../models/JobMatch";
import type { Profile } from "../models/Profile";
import type { ProfileInput } from "../models/ProfileInput";
import type { ReferenceLetter } from "../models/ReferenceLetter";
import type { CancelablePromise } from "../core/CancelablePromise";
import { OpenAPI } from "../core/OpenAPI";
import { request as __request } from "../core/request";
export class DefaultService {
	/**
	 * Get current user's profile
	 * @returns Profile Successful response
	 * @throws ApiError
	 */
	public static getProfile(): CancelablePromise<Profile> {
		return __request(OpenAPI, {
			method: "GET",
			url: "/profile",
		});
	}
	/**
	 * Update profile
	 * @param requestBody
	 * @returns Profile Updated successfully
	 * @throws ApiError
	 */
	public static updateProfile(
		requestBody: ProfileInput,
	): CancelablePromise<Profile> {
		return __request(OpenAPI, {
			method: "PUT",
			url: "/profile",
			body: requestBody,
			mediaType: "application/json",
		});
	}
	/**
	 * Upload a reference letter
	 * @param formData
	 * @returns ReferenceLetter Uploaded successfully
	 * @throws ApiError
	 */
	public static uploadReferenceLetter(formData: {
		file?: Blob;
	}): CancelablePromise<ReferenceLetter> {
		return __request(OpenAPI, {
			method: "POST",
			url: "/reference-letters",
			formData: formData,
			mediaType: "multipart/form-data",
		});
	}
	/**
	 * Generate profile from uploaded letters
	 * @returns Profile Generation started/completed
	 * @throws ApiError
	 */
	public static generateProfile(): CancelablePromise<Profile> {
		return __request(OpenAPI, {
			method: "POST",
			url: "/profile/generate",
		});
	}
	/**
	 * Create a tailored profile for a job description
	 * @param requestBody
	 * @returns JobMatch Tailored profile created
	 * @throws ApiError
	 */
	public static tailorProfile(requestBody: {
		jobDescription?: string;
	}): CancelablePromise<JobMatch> {
		return __request(OpenAPI, {
			method: "POST",
			url: "/profile/tailor",
			body: requestBody,
			mediaType: "application/json",
		});
	}
	/**
	 * Download CV as PDF
	 * @param profileId
	 * @returns binary PDF File
	 * @throws ApiError
	 */
	public static downloadCv(profileId: string): CancelablePromise<Blob> {
		return __request(OpenAPI, {
			method: "GET",
			url: "/profile/{profileId}/cv",
			path: {
				profileId: profileId,
			},
		});
	}
}
