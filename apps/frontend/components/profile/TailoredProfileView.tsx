"use client";

import type { Profile } from "../../lib/api/generated/models/Profile";
import type { JobMatch } from "../../lib/api/generated/models/JobMatch";
import WorkExperienceCard from "./WorkExperienceCard";
import SkillsSection from "./SkillsSection";
import MatchScore from "./MatchScore";

interface TailoredProfileViewProps {
	profile: Profile;
	jobMatch: JobMatch;
}

export default function TailoredProfileView({
	profile,
	jobMatch,
}: TailoredProfileViewProps) {
	const matchScore = jobMatch.matchScore || 0;
	const tailoredSummary = jobMatch.tailoredSummary || profile.summary || "";

	return (
		<div className="max-w-4xl mx-auto bg-gray-50 min-h-screen">
			<div className="bg-white shadow-lg rounded-lg p-8">
				{/* Match Score */}
				{matchScore > 0 && <MatchScore score={matchScore} />}

				{/* Header with Tailored Summary */}
				<div className="border-b border-gray-200 pb-6 mb-6">
					<div className="flex items-center justify-between mb-2">
						<h1 className="text-3xl font-bold text-gray-900">
							Tailored Profile
						</h1>
						<span className="px-3 py-1 bg-indigo-100 text-indigo-800 rounded-full text-sm font-medium">
							Tailored
						</span>
					</div>
					{tailoredSummary && (
						<div className="mt-4">
							<p className="text-gray-700 text-lg leading-relaxed">
								{tailoredSummary}
							</p>
							{jobMatch.tailoredSummary &&
								jobMatch.tailoredSummary !== profile.summary && (
									<p className="mt-3 text-sm text-indigo-600 italic">
										✓ Summary tailored to highlight relevant experience
									</p>
								)}
						</div>
					)}
				</div>

				{/* Skills Section */}
				{profile.skills && profile.skills.length > 0 && (
					<SkillsSection skills={profile.skills} />
				)}

				{/* Work Experience Section */}
				{profile.workExperiences && profile.workExperiences.length > 0 && (
					<div className="mb-6">
						<h2 className="text-2xl font-bold text-gray-900 mb-4">
							Work Experience
						</h2>
						<div className="space-y-4">
							{profile.workExperiences.map((exp) => (
								<WorkExperienceCard key={exp.id} experience={exp} />
							))}
						</div>
					</div>
				)}

				{/* Explanation Note */}
				<div className="mt-6 p-4 bg-indigo-50 border-l-4 border-indigo-500 rounded">
					<p className="text-sm text-indigo-800">
						<strong>Note:</strong> This profile has been tailored to emphasize
						skills and experiences most relevant to the job description. The
						original profile content is preserved, with emphasis added to highly
						relevant areas.
					</p>
				</div>
			</div>
		</div>
	);
}
