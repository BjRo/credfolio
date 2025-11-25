"use client";

import {
	CredibilityHighlights,
	type CredibilityHighlight,
} from "./CredibilityHighlights";

export interface WorkExperience {
	id: string;
	companyName: string;
	role: string;
	startDate: string;
	endDate?: string;
	description: string;
	credibilityHighlights: CredibilityHighlight[];
}

interface WorkExperienceCardProps {
	experience: WorkExperience;
}

const formatDate = (dateString: string): string => {
	const date = new Date(dateString);
	return date.toLocaleDateString("en-US", {
		month: "short",
		year: "numeric",
	});
};

export const WorkExperienceCard = ({ experience }: WorkExperienceCardProps) => {
	const hasHighlights = experience.credibilityHighlights.length > 0;
	const highlightCount = experience.credibilityHighlights.length;

	return (
		<div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-shadow">
			<div className="flex flex-col md:flex-row md:items-start md:justify-between gap-2 mb-4">
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
						{experience.endDate ? formatDate(experience.endDate) : "Present"}
					</span>
					{hasHighlights && (
						<span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-emerald-100 text-emerald-700 text-xs font-medium">
							<svg
								className="w-3.5 h-3.5"
								fill="currentColor"
								viewBox="0 0 20 20"
								aria-hidden="true"
							>
								<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
							</svg>
							{highlightCount} endorsement{highlightCount > 1 ? "s" : ""}
						</span>
					)}
				</div>
			</div>

			<p className="text-gray-600 mb-4 leading-relaxed">
				{experience.description}
			</p>

			{hasHighlights && (
				<div className="mt-4 pt-4 border-t border-gray-100">
					<h4 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
						<svg
							className="w-4 h-4 text-emerald-500"
							fill="currentColor"
							viewBox="0 0 20 20"
							aria-hidden="true"
						>
							<path
								fillRule="evenodd"
								d="M18 13V5a2 2 0 00-2-2H4a2 2 0 00-2 2v8a2 2 0 002 2h3l3 3 3-3h3a2 2 0 002-2zM5 7a1 1 0 011-1h8a1 1 0 110 2H6a1 1 0 01-1-1zm1 3a1 1 0 100 2h3a1 1 0 100-2H6z"
								clipRule="evenodd"
							/>
						</svg>
						Credibility Highlights
					</h4>
					<CredibilityHighlights
						highlights={experience.credibilityHighlights}
					/>
				</div>
			)}
		</div>
	);
};
