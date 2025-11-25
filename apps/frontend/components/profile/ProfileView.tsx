"use client";

import { useEffect, useState } from "react";
import { getProfile } from "@/lib/api/profile";
import { WorkExperienceCard, type WorkExperience } from "./WorkExperienceCard";
import { SkillsSection } from "./SkillsSection";

interface Profile {
	id: string;
	summary: string;
	workExperiences: WorkExperience[];
	skills: string[];
}

export const ProfileView = () => {
	const [profile, setProfile] = useState<Profile | null>(null);
	const [isLoading, setIsLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const fetchProfile = async () => {
			try {
				setIsLoading(true);
				const data = await getProfile();
				setProfile(data);
				setError(null);
			} catch (err) {
				setError("Failed to load profile");
				console.error("Error fetching profile:", err);
			} finally {
				setIsLoading(false);
			}
		};

		fetchProfile();
	}, []);

	if (isLoading) {
		return (
			<div className="flex items-center justify-center min-h-[400px]">
				<div className="flex flex-col items-center gap-4">
					<div className="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent" />
					<p className="text-gray-500 animate-pulse">Loading your profile...</p>
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div className="flex items-center justify-center min-h-[400px]">
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
					<p className="text-red-600 font-medium">{error}</p>
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
		<div className="max-w-4xl mx-auto space-y-8">
			{/* Hero Section with Summary */}
			<section className="bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 rounded-2xl p-8 text-white shadow-xl">
				<div className="flex items-center gap-4 mb-6">
					<div className="w-20 h-20 rounded-full bg-white/20 backdrop-blur flex items-center justify-center">
						<svg
							className="w-10 h-10 text-white"
							fill="currentColor"
							viewBox="0 0 20 20"
							aria-hidden="true"
						>
							<path
								fillRule="evenodd"
								d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
								clipRule="evenodd"
							/>
						</svg>
					</div>
					<div>
						<h1 className="text-3xl font-bold">Professional Profile</h1>
						<p className="text-white/80">
							Your verified credibility at a glance
						</p>
					</div>
				</div>

				<div className="bg-white/10 backdrop-blur rounded-xl p-6">
					<h2 className="text-sm font-semibold uppercase tracking-wider text-white/70 mb-2">
						Summary
					</h2>
					{profile.summary ? (
						<p className="text-lg leading-relaxed">{profile.summary}</p>
					) : (
						<p className="text-white/60 italic">No summary available</p>
					)}
				</div>
			</section>

			{/* Skills Section */}
			<section className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
				<SkillsSection skills={profile.skills} />
			</section>

			{/* Work Experience Section */}
			<section>
				<div className="flex items-center gap-3 mb-6">
					<h2 className="text-2xl font-bold text-gray-900">Work Experience</h2>
					<span className="px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-sm font-medium">
						{profile.workExperiences.length} position
						{profile.workExperiences.length !== 1 ? "s" : ""}
					</span>
				</div>

				{profile.workExperiences.length > 0 ? (
					<div className="space-y-4">
						{profile.workExperiences.map((experience) => (
							<WorkExperienceCard key={experience.id} experience={experience} />
						))}
					</div>
				) : (
					<div className="bg-gray-50 rounded-xl p-8 text-center">
						<svg
							className="w-12 h-12 text-gray-400 mx-auto mb-4"
							fill="currentColor"
							viewBox="0 0 20 20"
							aria-hidden="true"
						>
							<path
								fillRule="evenodd"
								d="M6 6V5a3 3 0 013-3h2a3 3 0 013 3v1h2a2 2 0 012 2v3.57A22.952 22.952 0 0110 13a22.95 22.95 0 01-8-1.43V8a2 2 0 012-2h2zm2-1a1 1 0 011-1h2a1 1 0 011 1v1H8V5zm1 5a1 1 0 011-1h.01a1 1 0 110 2H10a1 1 0 01-1-1z"
								clipRule="evenodd"
							/>
							<path d="M2 13.692V16a2 2 0 002 2h12a2 2 0 002-2v-2.308A24.974 24.974 0 0110 15c-2.796 0-5.487-.46-8-1.308z" />
						</svg>
						<p className="text-gray-500">No work experience yet</p>
						<p className="text-gray-400 text-sm mt-1">
							Upload reference letters to generate your experience
						</p>
					</div>
				)}
			</section>

			{/* Total Credibility Stats */}
			{profile.workExperiences.length > 0 && (
				<section className="bg-gradient-to-r from-emerald-50 to-teal-50 rounded-xl p-6 border border-emerald-100">
					<h3 className="text-lg font-semibold text-emerald-800 mb-4 flex items-center gap-2">
						<svg
							className="w-5 h-5"
							fill="currentColor"
							viewBox="0 0 20 20"
							aria-hidden="true"
						>
							<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
						</svg>
						Credibility Overview
					</h3>
					<div className="grid grid-cols-3 gap-4 text-center">
						<div className="bg-white rounded-lg p-4 shadow-sm">
							<div className="text-3xl font-bold text-emerald-600">
								{profile.workExperiences.reduce(
									(sum, exp) => sum + exp.credibilityHighlights.length,
									0,
								)}
							</div>
							<div className="text-sm text-gray-600">Total Endorsements</div>
						</div>
						<div className="bg-white rounded-lg p-4 shadow-sm">
							<div className="text-3xl font-bold text-indigo-600">
								{profile.workExperiences.length}
							</div>
							<div className="text-sm text-gray-600">Verified Positions</div>
						</div>
						<div className="bg-white rounded-lg p-4 shadow-sm">
							<div className="text-3xl font-bold text-purple-600">
								{profile.skills.length}
							</div>
							<div className="text-sm text-gray-600">Skills</div>
						</div>
					</div>
				</section>
			)}
		</div>
	);
};
