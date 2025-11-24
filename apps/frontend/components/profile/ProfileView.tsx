"use client";

import type { Profile } from "../../lib/api/generated/models/Profile";
import WorkExperienceCard from "./WorkExperienceCard";
import SkillsSection from "./SkillsSection";

export default function ProfileView({ profile }: { profile: Profile }) {
	return (
		<div className="max-w-4xl mx-auto bg-gray-50 min-h-screen">
			<div className="bg-white shadow-lg rounded-lg p-8">
				{/* Header */}
				<div className="border-b border-gray-200 pb-6 mb-6">
					<h1 className="text-3xl font-bold text-gray-900 mb-2">
						Professional Profile
					</h1>
					{profile.summary && (
						<p className="text-gray-700 text-lg leading-relaxed">
							{profile.summary}
						</p>
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

				{/* Empty State */}
				{(!profile.workExperiences ||
					profile.workExperiences.length === 0) &&
					(!profile.skills || profile.skills.length === 0) && (
						<div className="text-center py-12 text-gray-500">
							<p>No profile data available yet.</p>
							<p className="mt-2 text-sm">
								Upload reference letters and generate your profile to get
								started.
							</p>
						</div>
					)}
			</div>
		</div>
	);
}

