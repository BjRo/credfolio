"use client";

interface SkillsSectionProps {
	skills: string[];
}

export const SkillsSection = ({ skills }: SkillsSectionProps) => {
	// Remove duplicates and sort alphabetically
	const uniqueSkills = [...new Set(skills)].sort((a, b) =>
		a.toLowerCase().localeCompare(b.toLowerCase()),
	);

	if (uniqueSkills.length === 0) {
		return (
			<div className="text-gray-500 italic text-sm">No skills added yet</div>
		);
	}

	return (
		<div>
			<div className="flex items-center justify-between mb-4">
				<h3 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
					<svg
						className="w-5 h-5 text-indigo-500"
						fill="currentColor"
						viewBox="0 0 20 20"
						aria-hidden="true"
					>
						<path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838L7.667 9.088l1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zM9.3 16.573A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z" />
					</svg>
					Skills
				</h3>
				<span className="text-sm text-gray-500">
					{uniqueSkills.length} skill{uniqueSkills.length !== 1 ? "s" : ""}
				</span>
			</div>
			<ul className="flex flex-wrap gap-2">
				{uniqueSkills.map((skill) => (
					<li
						key={skill}
						className="px-3 py-1.5 bg-gradient-to-r from-indigo-50 to-purple-50 text-indigo-700 rounded-full text-sm font-medium border border-indigo-100 hover:from-indigo-100 hover:to-purple-100 transition-colors"
					>
						{skill}
					</li>
				))}
			</ul>
		</div>
	);
};
