export default function SkillsList({ skills, highlightIDs }: { skills: any[]; highlightIDs?: string[] }) {
	if (!skills || skills.length === 0) return null;
	return (
		<div className="flex flex-wrap gap-2">
			{skills.map((skill) => {
				const isHighlighted = highlightIDs?.includes(skill.id);
				return (
					<span
						key={skill.id}
						className={`px-2 py-1 rounded text-xs font-medium transition-colors ${
							isHighlighted
								? 'bg-green-200 text-green-800 dark:bg-green-900 dark:text-green-200 ring-2 ring-green-500'
								: 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200'
						}`}
					>
						{skill.name}
					</span>
				);
			})}
		</div>
	);
}
