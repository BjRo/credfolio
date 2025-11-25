"use client";

interface MatchScoreProps {
	score: number;
	summary?: string;
}

export const MatchScore = ({ score, summary }: MatchScoreProps) => {
	const percentage = Math.round(score * 100);

	const getMatchLevel = () => {
		if (score >= 0.7) {
			return {
				label: "Strong Match",
				color: "text-emerald-600",
				bgColor: "bg-emerald-500",
				lightBg: "from-emerald-50 to-teal-50",
				borderColor: "border-emerald-200",
			};
		}
		if (score >= 0.4) {
			return {
				label: "Moderate Match",
				color: "text-amber-600",
				bgColor: "bg-amber-500",
				lightBg: "from-amber-50 to-yellow-50",
				borderColor: "border-amber-200",
			};
		}
		return {
			label: "Limited Match",
			color: "text-red-600",
			bgColor: "bg-red-500",
			lightBg: "from-red-50 to-orange-50",
			borderColor: "border-red-200",
		};
	};

	const matchLevel = getMatchLevel();

	return (
		<div
			className={`rounded-xl p-6 bg-gradient-to-r ${matchLevel.lightBg} border ${matchLevel.borderColor}`}
		>
			<div className="flex items-center justify-between mb-4">
				<div>
					<h3 className="text-lg font-semibold text-gray-900">Match Score</h3>
					<p className={`text-sm font-medium ${matchLevel.color}`}>
						{matchLevel.label}
					</p>
				</div>
				<div className="relative">
					<div className="w-20 h-20 relative">
						<svg
							className="w-full h-full -rotate-90"
							viewBox="0 0 36 36"
							aria-hidden="true"
						>
							<path
								d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
								fill="none"
								stroke="#e5e7eb"
								strokeWidth="3"
							/>
							<path
								d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
								fill="none"
								stroke="currentColor"
								className={matchLevel.color}
								strokeWidth="3"
								strokeDasharray={`${percentage}, 100`}
								strokeLinecap="round"
							/>
						</svg>
						<div className="absolute inset-0 flex items-center justify-center">
							<span className="text-2xl font-bold text-gray-900">
								{percentage}%
							</span>
						</div>
					</div>
				</div>
			</div>

			{/* Progress bar */}
			<div className="h-2 bg-gray-200 rounded-full overflow-hidden">
				<div
					className={`h-full ${matchLevel.bgColor} rounded-full transition-all duration-500`}
					style={{ width: `${percentage}%` }}
				/>
			</div>

			{summary && (
				<p className="mt-4 text-gray-600 text-sm leading-relaxed">{summary}</p>
			)}
		</div>
	);
};
