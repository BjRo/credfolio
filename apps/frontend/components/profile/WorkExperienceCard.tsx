"use client";

import type { WorkExperience } from "../../lib/api/generated/models/WorkExperience";
import CredibilityHighlights from "./CredibilityHighlights";

function formatDate(dateStr?: string): string {
	if (!dateStr) return "";
	try {
		const date = new Date(dateStr);
		return date.toLocaleDateString("en-US", {
			year: "numeric",
			month: "short",
		});
	} catch {
		return dateStr;
	}
}

export default function WorkExperienceCard({
	experience,
}: {
	experience: WorkExperience;
}) {
	const startDate = formatDate(experience.startDate);
	const endDate = experience.endDate
		? formatDate(experience.endDate)
		: "Present";

	return (
		<div className="bg-white rounded-lg shadow-md p-6 mb-6 border-l-4 border-indigo-500">
			<div className="flex justify-between items-start mb-2">
				<div>
					<h3 className="text-xl font-bold text-gray-900">{experience.role}</h3>
					<p className="text-lg text-indigo-600 font-medium">
						{experience.companyName}
					</p>
				</div>
				<div className="text-right text-sm text-gray-600">
					{startDate} - {endDate}
				</div>
			</div>
			{experience.description && (
				<p className="text-gray-700 mt-3 leading-relaxed">
					{experience.description}
				</p>
			)}
			<CredibilityHighlights highlights={experience.credibilityHighlights} />
		</div>
	);
}
