"use client";

import { useEffect, useState } from "react";
import { getProfile, type Profile } from "@/lib/api/profile";
import { ProfileView } from "@/components/profile/ProfileView";
import { DownloadCVButton } from "@/components/profile/DownloadCVButton";
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
			<div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-indigo-50 flex items-center justify-center">
				<div className="flex flex-col items-center gap-4">
					<div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent" />
					<p className="text-gray-500 animate-pulse">Loading your profile...</p>
				</div>
			</div>
		);
	}

	if (error === "noProfile") {
		return (
			<div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-indigo-50">
				<div className="max-w-2xl mx-auto px-4 py-24 text-center">
					<div className="w-24 h-24 mx-auto mb-8 rounded-full bg-gradient-to-br from-indigo-100 to-purple-100 flex items-center justify-center shadow-lg">
						<svg
							className="w-12 h-12 text-indigo-600"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
							/>
						</svg>
					</div>
					<h1 className="text-4xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-4">
						No Profile Yet
					</h1>
					<p className="text-lg text-gray-600 mb-8 max-w-md mx-auto">
						Upload your reference letters to generate a professional profile
						with AI-powered credibility insights.
					</p>
					<Link
						href="/profile/generate"
						className="inline-flex items-center gap-2 px-8 py-4 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-xl hover:from-indigo-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-105 shadow-lg hover:shadow-xl"
					>
						<svg
							className="w-5 h-5"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
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
			<div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-indigo-50 flex items-center justify-center">
				<div className="text-center">
					<div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-100 mb-4">
						<svg
							className="w-8 h-8 text-red-500"
							fill="currentColor"
							viewBox="0 0 20 20"
							aria-hidden="true"
						>
							<path
								fillRule="evenodd"
								d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
								clipRule="evenodd"
							/>
						</svg>
					</div>
					<p className="text-red-600 text-lg font-medium">{error}</p>
					<button
						type="button"
						onClick={() => window.location.reload()}
						className="mt-4 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
					>
						Try Again
					</button>
				</div>
			</div>
		);
	}

	if (!profile) {
		return null;
	}

	return (
		<div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-indigo-50">
			{/* Navigation Bar */}
			<nav className="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-gray-200/50">
				<div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
					<h1 className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
						Credfolio
					</h1>
					<div className="flex gap-3">
						<DownloadCVButton profileId={profile.id} />
						<Link
							href="/profile/generate"
							className="px-4 py-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors text-sm font-medium"
						>
							Add Letters
						</Link>
						<Link
							href="/profile/tailor"
							className="px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg hover:from-purple-700 hover:to-pink-700 transition-all text-sm font-medium shadow-sm"
						>
							Tailor to Job
						</Link>
					</div>
				</div>
			</nav>

			{/* Main Content */}
			<main className="max-w-4xl mx-auto px-4 py-8">
				<ProfileView />
			</main>
		</div>
	);
}
