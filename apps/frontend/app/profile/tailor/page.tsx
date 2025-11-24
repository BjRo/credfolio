"use client";

import { useState } from "react";
import { tailorProfile } from "../../../lib/api/profile";
import { getProfile } from "../../../lib/api/profile";
import type { Profile } from "../../../lib/api/generated/models/Profile";
import type { JobMatch } from "../../../lib/api/generated/models/JobMatch";
import JobDescriptionInput from "../../../components/profile/JobDescriptionInput";
import TailoredProfileView from "../../../components/profile/TailoredProfileView";
import ProfileView from "../../../components/profile/ProfileView";
import DownloadCVButton from "../../../components/profile/DownloadCVButton";

export default function TailorProfilePage() {
	const [profile, setProfile] = useState<Profile | null>(null);
	const [jobMatch, setJobMatch] = useState<JobMatch | null>(null);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [showTailored, setShowTailored] = useState(false);

	const loadProfile = async () => {
		try {
			setError(null);
			const data = await getProfile();
			setProfile(data);
		} catch (err) {
			console.error("Failed to load profile:", err);
			setError("Failed to load profile. Please try again.");
		}
	};

	const handleTailor = async (jobDescription: string) => {
		try {
			setLoading(true);
			setError(null);
			const match = await tailorProfile(jobDescription);
			setJobMatch(match);
			setShowTailored(true);

			// Reload profile to ensure we have the latest data
			await loadProfile();
		} catch (err) {
			console.error("Failed to tailor profile:", err);
			setError("Failed to tailor profile. Please try again.");
		} finally {
			setLoading(false);
		}
	};

	// Load profile on mount
	if (!profile && !loading && !error) {
		loadProfile();
	}

	return (
		<div className="max-w-4xl mx-auto py-8 px-4">
			<h1 className="text-3xl font-bold text-gray-900 mb-6">
				Tailor Your Profile
			</h1>

			{/* Error Message */}
			{error && (
				<div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
					{error}
					<button
						type="button"
						onClick={() => {
							setError(null);
							loadProfile();
						}}
						className="ml-4 text-red-600 underline"
					>
						Retry
					</button>
				</div>
			)}

			{/* Job Description Input */}
			<div className="bg-white rounded-lg shadow-md p-6 mb-6">
				<JobDescriptionInput onSubmit={handleTailor} isLoading={loading} />
			</div>

			{/* Toggle between Standard and Tailored View */}
			{profile && jobMatch && (
				<div className="mb-6 flex gap-4">
					<button
						type="button"
						onClick={() => setShowTailored(false)}
						className={`px-4 py-2 rounded-lg transition-colors ${
							!showTailored
								? "bg-indigo-600 text-white"
								: "bg-gray-200 text-gray-700 hover:bg-gray-300"
						}`}
					>
						Standard Profile
					</button>
					<button
						type="button"
						onClick={() => setShowTailored(true)}
						className={`px-4 py-2 rounded-lg transition-colors ${
							showTailored
								? "bg-indigo-600 text-white"
								: "bg-gray-200 text-gray-700 hover:bg-gray-300"
						}`}
					>
						Tailored Profile
					</button>
				</div>
			)}

			{/* Download Button */}
			{profile?.id && (
				<div className="mb-4 flex justify-end">
					<DownloadCVButton
						profileId={profile.id}
						jobMatchId={showTailored && jobMatch?.id ? jobMatch.id : undefined}
					/>
				</div>
			)}

			{/* Profile Display */}
			{profile &&
				(showTailored && jobMatch ? (
					<TailoredProfileView profile={profile} jobMatch={jobMatch} />
				) : (
					<ProfileView profile={profile} />
				))}

			{/* Loading State */}
			{loading && !profile && (
				<div className="text-center py-12">
					<div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto" />
					<p className="mt-4 text-gray-600">Tailoring profile...</p>
				</div>
			)}

			{/* No Profile State */}
			{!profile && !loading && !error && (
				<div className="text-center py-12 text-gray-500">
					<p>No profile found. Please generate your profile first.</p>
				</div>
			)}
		</div>
	);
}
