"use client";

import { MatchScore } from "./MatchScore";
import { SkillsSection } from "./SkillsSection";

interface TailoredExperience {
	id: string;
	companyName: string;
	role: string;
	startDate: string;
	endDate?: string;
	description: string;
	relevanceScore: number;
	highlightReason?: string;
}

interface TailoredProfileData {
	id: string;
	matchScore: number;
	matchSummary: string;
	tailoredExperiences: TailoredExperience[];
	relevantSkills: string[];
}

interface TailoredProfileViewProps {
	tailoredProfile: TailoredProfileData;
}

const formatDate = (dateString: string): string => {
	const date = new Date(dateString);
	return date.toLocaleDateString("en-US", {
		month: "short",
		year: "numeric",
	});
};

export const TailoredProfileView = ({
	tailoredProfile,
}: TailoredProfileViewProps) => {
	return (
		<div className="space-y-8">
			{/* Match Score Section */}
			<MatchScore
				score={tailoredProfile.matchScore}
				summary={tailoredProfile.matchSummary}
			/>

			{/* Relevant Skills */}
			<section className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
				<SkillsSection skills={tailoredProfile.relevantSkills} />
			</section>

			{/* Tailored Experiences */}
			<section>
				<div className="flex items-center gap-3 mb-6">
					<h2 className="text-2xl font-bold text-gray-900">
						Ranked Experience
					</h2>
					<span className="px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-sm font-medium">
						By Relevance
					</span>
				</div>

				<div className="space-y-4">
					{tailoredProfile.tailoredExperiences.map((experience, index) => (
						<TailoredExperienceCard
							key={experience.id}
							experience={experience}
							rank={index + 1}
						/>
					))}
				</div>
			</section>
		</div>
	);
};

interface TailoredExperienceCardProps {
	experience: TailoredExperience;
	rank: number;
}

const TailoredExperienceCard = ({
	experience,
	rank,
}: TailoredExperienceCardProps) => {
	const relevancePercentage = Math.round(experience.relevanceScore * 100);
	const isHighRelevance = experience.relevanceScore >= 0.7;
	const isMediumRelevance =
		experience.relevanceScore >= 0.4 && experience.relevanceScore < 0.7;

	const relevanceColor = isHighRelevance
		? "text-emerald-600 bg-emerald-50 border-emerald-200"
		: isMediumRelevance
			? "text-amber-600 bg-amber-50 border-amber-200"
			: "text-gray-500 bg-gray-50 border-gray-200";

	return (
		<article
			className={`bg-white rounded-xl shadow-sm border p-6 transition-all ${
				isHighRelevance
					? "border-emerald-200 ring-1 ring-emerald-100"
					: "border-gray-100"
			}`}
		>
			<div className="flex items-start gap-4">
				{/* Rank Badge */}
				<div
					className={`w-10 h-10 rounded-full flex items-center justify-center font-bold ${
						rank === 1
							? "bg-gradient-to-br from-amber-400 to-amber-500 text-white"
							: rank === 2
								? "bg-gradient-to-br from-gray-300 to-gray-400 text-white"
								: rank === 3
									? "bg-gradient-to-br from-amber-600 to-amber-700 text-white"
									: "bg-gray-100 text-gray-500"
					}`}
				>
					#{rank}
				</div>

				<div className="flex-1">
					<div className="flex flex-col md:flex-row md:items-start md:justify-between gap-2 mb-2">
						<div>
							<h3 className="text-xl font-semibold text-gray-900">
								{experience.companyName}
							</h3>
							<p className="text-lg text-indigo-600 font-medium">
								{experience.role}
							</p>
						</div>
						<div className="flex items-center gap-2">
							<span className="text-sm text-gray-500">
								{formatDate(experience.startDate)} -{" "}
								{experience.endDate
									? formatDate(experience.endDate)
									: "Present"}
							</span>
							<span
								className={`px-3 py-1 rounded-full text-sm font-bold border ${relevanceColor}`}
							>
								{relevancePercentage}%
							</span>
						</div>
					</div>

					<p className="text-gray-600 mb-3">{experience.description}</p>

					{experience.highlightReason && (
						<div className="mt-3 p-3 bg-indigo-50 rounded-lg border border-indigo-100">
							<div className="flex items-start gap-2">
								<svg
									className="w-5 h-5 text-indigo-500 mt-0.5 flex-shrink-0"
									fill="currentColor"
									viewBox="0 0 20 20"
									aria-hidden="true"
								>
									<path
										fillRule="evenodd"
										d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
										clipRule="evenodd"
									/>
								</svg>
								<p className="text-sm text-indigo-700">
									{experience.highlightReason}
								</p>
							</div>
						</div>
					)}
				</div>
			</div>
		</article>
	);
};
