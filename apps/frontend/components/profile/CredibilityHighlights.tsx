"use client";

export interface CredibilityHighlight {
	quote: string;
	sentiment: "POSITIVE" | "NEUTRAL";
}

interface CredibilityHighlightsProps {
	highlights: CredibilityHighlight[];
}

export const CredibilityHighlights = ({
	highlights,
}: CredibilityHighlightsProps) => {
	if (highlights.length === 0) {
		return (
			<div className="text-gray-500 italic text-sm">
				No credibility highlights yet
			</div>
		);
	}

	return (
		<div className="space-y-3">
			{highlights.map((highlight) => (
				<blockquote
					key={highlight.quote}
					data-sentiment={highlight.sentiment}
					className={`border-l-4 pl-4 py-2 ${
						highlight.sentiment === "POSITIVE"
							? "border-emerald-500 bg-emerald-50"
							: "border-slate-400 bg-slate-50"
					}`}
				>
					<p className="text-gray-700 italic">
						&ldquo;{highlight.quote}&rdquo;
					</p>
					<div className="mt-1 flex items-center gap-2">
						{highlight.sentiment === "POSITIVE" ? (
							<span className="inline-flex items-center gap-1 text-emerald-600 text-xs font-medium">
								<svg
									className="w-3.5 h-3.5"
									fill="currentColor"
									viewBox="0 0 20 20"
									aria-hidden="true"
								>
									<path
										fillRule="evenodd"
										d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
										clipRule="evenodd"
									/>
								</svg>
								Positive
							</span>
						) : (
							<span className="text-slate-500 text-xs font-medium">
								Neutral
							</span>
						)}
					</div>
				</blockquote>
			))}
		</div>
	);
};
