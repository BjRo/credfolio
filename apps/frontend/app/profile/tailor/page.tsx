"use client";

import { useState } from "react";
import Link from "next/link";
import { JobDescriptionInput } from "@/components/profile/JobDescriptionInput";
import { TailoredProfileView } from "@/components/profile/TailoredProfileView";
import { DownloadCVButton } from "@/components/profile/DownloadCVButton";
import {
	tailorProfile,
	type TailoredProfile,
	getProfile,
} from "@/lib/api/profile";

export default function TailorProfilePage() {
	const [jobDescription, setJobDescription] = useState("");
	const [tailoredProfile, setTailoredProfile] =
		useState<TailoredProfile | null>(null);
	const [profileId, setProfileId] = useState<string | null>(null);
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleTailor = async () => {
		if (!jobDescription.trim()) {
			setError("Please enter a job description");
			return;
		}

		setIsLoading(true);
		setError(null);

		try {
			// Get profile ID first
			const profile = await getProfile();
			setProfileId(profile.id);

			const result = await tailorProfile(jobDescription);
			setTailoredProfile(result);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Failed to tailor profile");
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-purple-50">
			{/* Navigation */}
			<nav className="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-gray-200/50">
				<div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
					<Link
						href="/profile"
						className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent"
					>
						Credfolio
					</Link>
					<div className="flex gap-3">
						<Link
							href="/profile"
							className="px-4 py-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors text-sm font-medium"
						>
							View Profile
						</Link>
						<Link
							href="/profile/generate"
							className="px-4 py-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors text-sm font-medium"
						>
							Add Letters
						</Link>
					</div>
				</div>
			</nav>

			<main className="max-w-4xl mx-auto px-4 py-8">
				{/* Header */}
				<div className="mb-8">
					<h1 className="text-3xl font-bold text-gray-900 mb-2">
						Tailor Your Profile
					</h1>
					<p className="text-gray-600">
						Paste a job description to see how your experience matches and get
						AI-powered recommendations.
					</p>
				</div>

				{/* Input Section */}
				{!tailoredProfile && (
					<div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 mb-8">
						<JobDescriptionInput
							value={jobDescription}
							onChange={setJobDescription}
							onTailor={handleTailor}
							isLoading={isLoading}
							error={error}
						/>
					</div>
				)}

				{/* Results Section */}
				{tailoredProfile && (
					<>
						<div className="mb-6 flex items-center justify-between">
							<h2 className="text-xl font-semibold text-gray-900">
								Tailored Results
							</h2>
							<div className="flex items-center gap-3">
								{profileId && (
									<DownloadCVButton
										profileId={profileId}
										tailoredId={tailoredProfile.id}
									/>
								)}
								<button
									type="button"
									onClick={() => {
										setTailoredProfile(null);
										setJobDescription("");
									}}
									className="text-indigo-600 hover:text-indigo-700 text-sm font-medium"
								>
									← Try Another Job
								</button>
							</div>
						</div>
						<TailoredProfileView tailoredProfile={tailoredProfile} />
					</>
				)}
			</main>
		</div>
	);
}
