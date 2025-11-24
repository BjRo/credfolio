"use client";

export default function SkillsSection({
	skills,
}: {
	skills?: Array<string>;
}) {
	if (!skills || skills.length === 0) {
		return null;
	}

	return (
		<div className="mb-6">
			<h3 className="text-xl font-bold text-gray-900 mb-4">Skills</h3>
			<div className="flex flex-wrap gap-2">
				{skills.map((skill) => (
					<span
						key={skill}
						className="px-3 py-1 bg-indigo-100 text-indigo-800 rounded-full text-sm font-medium"
					>
						{skill}
					</span>
				))}
			</div>
		</div>
	);
}
