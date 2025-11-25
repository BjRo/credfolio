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
							disabled={isLoading}
						/>

						{error && (
							<div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
								<p className="text-red-600 text-sm">{error}</p>
							</div>
						)}

						<div className="mt-6 flex justify-end">
							<button
								type="button"
								onClick={handleTailor}
								disabled={isLoading || !jobDescription.trim()}
								className={`px-6 py-3 rounded-xl font-semibold transition-all ${
									isLoading || !jobDescription.trim()
										? "bg-gray-200 text-gray-400 cursor-not-allowed"
										: "bg-gradient-to-r from-purple-600 to-pink-600 text-white hover:from-purple-700 hover:to-pink-700 shadow-md hover:shadow-lg transform hover:scale-105"
								}`}
							>
								{isLoading ? (
									<span className="flex items-center gap-2">
										<svg
											className="animate-spin w-5 h-5"
											viewBox="0 0 24 24"
											aria-hidden="true"
										>
											<circle
												className="opacity-25"
												cx="12"
												cy="12"
												r="10"
												stroke="currentColor"
												strokeWidth="4"
												fill="none"
											/>
											<path
												className="opacity-75"
												fill="currentColor"
												d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
											/>
										</svg>
										Analyzing...
									</span>
								) : (
									<span className="flex items-center gap-2">
										<svg
											className="w-5 h-5"
											fill="currentColor"
											viewBox="0 0 20 20"
											aria-hidden="true"
										>
											<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
										</svg>
										Tailor Profile
									</span>
								)}
							</button>
						</div>
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
