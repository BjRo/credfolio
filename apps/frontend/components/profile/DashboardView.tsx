"use client";

import { useState } from 'react';
import type { UserProfile, TailoringResult } from '../../types';
import ProfileHeader from './ProfileHeader';
import ExperienceList from './ExperienceList';
import CredibilitySection from './CredibilitySection';
import JobDescriptionInput from './JobDescriptionInput';

export default function DashboardView({ profile, userID }: { profile: UserProfile; userID: string }) {
	const [tailoring, setTailoring] = useState<TailoringResult | null>(null);

	const handleTailored = (res: TailoringResult) => {
		setTailoring(res);
	};

	return (
		<div className="max-w-4xl mx-auto">
			<ProfileHeader profile={profile} />
			<CredibilitySection companies={profile.companies} />

			<JobDescriptionInput userID={userID} onTailored={handleTailored} />

			{tailoring && (
				<div className="mb-8 p-6 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-100 dark:border-green-900">
					<div className="flex items-center justify-between mb-4">
						<h3 className="text-xl font-bold text-green-800 dark:text-green-300">Tailored Match</h3>
						<span className="text-3xl font-bold text-green-600 dark:text-green-400">{tailoring.match_score}%</span>
					</div>
					<p className="text-gray-700 dark:text-gray-300 italic text-lg leading-relaxed">"{tailoring.summary_highlights}"</p>
				</div>
			)}

			<div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-8 border border-gray-100 dark:border-gray-700">
				<h2 className="text-2xl font-bold mb-8 text-gray-900 dark:text-white border-b border-gray-100 dark:border-gray-700 pb-4">Professional Experience</h2>
				<ExperienceList companies={profile.companies} highlightSkillIDs={tailoring?.relevant_skill_ids} />
			</div>
		</div>
	);
}

