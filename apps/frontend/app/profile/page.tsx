"use client";

import { useEffect, useState } from "react";
import { getProfile } from "../../lib/api/profile";
import type { Profile } from "../../lib/api/generated/models/Profile";
import ProfileView from "../../components/profile/ProfileView";
import DownloadCVButton from "../../components/profile/DownloadCVButton";
import { getErrorMessage } from "../../lib/utils/errorMessages";

export default function ProfilePage() {
	const [profile, setProfile] = useState<Profile | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const loadProfile = async () => {
			try {
				setLoading(true);
				setError(null);
				const data = await getProfile();
				setProfile(data);
			} catch (err) {
				console.error("Failed to load profile:", err);
				setError(getErrorMessage(err));
			} finally {
				setLoading(false);
			}
		};

		loadProfile();
	}, []);

	if (loading) {
		return (
			<div className="max-w-4xl mx-auto bg-gray-50 min-h-screen flex items-center justify-center">
				<div className="text-center">
					<div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto" />
					<p className="mt-4 text-gray-600">Loading profile...</p>
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="max-w-4xl mx-auto bg-gray-50 min-h-screen flex items-center justify-center">
				<div className="text-center bg-white p-8 rounded-lg shadow-md">
					<p className="text-red-600 mb-4">{error}</p>
					<button
						type="button"
						onClick={() => window.location.reload()}
						className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700"
					>
						Retry
					</button>
				</div>
			</div>
		);
	}

	if (!profile) {
		return (
			<div className="max-w-4xl mx-auto bg-gray-50 min-h-screen flex items-center justify-center">
				<div className="text-center bg-white p-8 rounded-lg shadow-md">
					<p className="text-gray-600 mb-4">No profile found.</p>
					<a
						href="/profile/generate"
						className="text-indigo-600 hover:text-indigo-700 underline"
					>
						Generate your profile
					</a>
				</div>
			</div>
		);
	}

	return (
		<div className="bg-gray-100 py-8">
			<div className="max-w-4xl mx-auto px-4 mb-4 flex justify-end">
				{profile.id && <DownloadCVButton profileId={profile.id} />}
			</div>
			<ProfileView profile={profile} />
		</div>
	);
}
