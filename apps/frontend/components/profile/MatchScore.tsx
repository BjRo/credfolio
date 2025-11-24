"use client";

interface MatchScoreProps {
	score: number; // 0.0 to 1.0
}

export default function MatchScore({ score }: MatchScoreProps) {
	const percentage = Math.round(score * 100);
	const colorClass =
		percentage >= 80
			? "bg-green-500"
			: percentage >= 60
				? "bg-yellow-500"
				: "bg-red-500";

	const barColorClass =
		percentage >= 80
			? "bg-green-600"
			: percentage >= 60
				? "bg-yellow-600"
				: "bg-red-600";

	return (
		<div className="mb-6">
			<div className="flex items-center justify-between mb-2">
				<h3 className="text-xl font-bold text-gray-900">Match Score</h3>
				<span
					className={`px-3 py-1 rounded-full text-white font-semibold ${colorClass}`}
				>
					{percentage}%
				</span>
			</div>
			<div className="w-full bg-gray-200 rounded-full h-4 overflow-hidden">
				<div
					className={`h-full ${barColorClass} transition-all duration-500 ease-out`}
					style={{ width: `${percentage}%` }}
				/>
			</div>
			<p className="mt-2 text-sm text-gray-600">
				{percentage >= 80
					? "Excellent match! Your profile aligns well with this job."
					: percentage >= 60
						? "Good match. Your profile has strong relevance to this position."
						: "Moderate match. Consider highlighting more relevant experience."}
			</p>
		</div>
	);
}
