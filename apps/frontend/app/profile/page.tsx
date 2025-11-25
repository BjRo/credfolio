"use client";

import { useEffect, useState } from "react";
import { getProfile, type Profile } from "@/lib/api/profile";
import { ProfileEditor } from "@/components/profile/ProfileEditor";
import Link from "next/link";

export default function ProfilePage() {
	const [profile, setProfile] = useState<Profile | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const fetchProfile = async () => {
			try {
				const data = await getProfile();
				setProfile(data);
			} catch (err) {
				if (err instanceof Error && err.message === "Profile not found") {
					setError("noProfile");
				} else {
					setError("Failed to load profile");
				}
			} finally {
				setLoading(false);
			}
		};

		fetchProfile();
	}, []);

	if (loading) {
		return (
			<div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 flex items-center justify-center">
				<div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
			</div>
		);
	}

	if (error === "noProfile") {
		return (
			<div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
				<div className="max-w-2xl mx-auto px-4 py-24 text-center">
					<div className="w-24 h-24 mx-auto mb-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
						<svg
							className="w-12 h-12 text-blue-600 dark:text-blue-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<title>Profile</title>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
							/>
						</svg>
					</div>
					<h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
						No Profile Yet
					</h1>
					<p className="text-lg text-gray-600 dark:text-gray-400 mb-8">
						Upload your reference letters to generate a professional profile
						with AI-powered insights.
					</p>
					<Link
						href="/profile/generate"
						className="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-105"
					>
						<svg
							className="w-5 h-5"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<title>Generate</title>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M13 10V3L4 14h7v7l9-11h-7z"
							/>
						</svg>
						Generate Your Profile
					</Link>
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 flex items-center justify-center">
				<div className="text-center">
					<p className="text-red-600 dark:text-red-400 text-lg">{error}</p>
				</div>
			</div>
		);
	}

	if (!profile) {
		return null;
	}

	return (
		<div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
			<div className="max-w-4xl mx-auto px-4 py-12">
				<div className="flex items-center justify-between mb-8">
					<h1 className="text-3xl font-bold text-gray-900 dark:text-white">
						Your Profile
					</h1>
					<div className="flex gap-4">
						<Link
							href="/profile/tailor"
							className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
						>
							Tailor to Job
						</Link>
						<Link
							href="/profile/generate"
							className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
						>
							Add More Letters
						</Link>
					</div>
				</div>

				<ProfileEditor
					profile={profile}
					onUpdate={setProfile}
					onError={(err) => setError(err)}
				/>
			</div>
		</div>
	);
}
