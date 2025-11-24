"use client";

import type { CredibilityHighlight } from "../../lib/api/generated/models/CredibilityHighlight";

export default function CredibilityHighlights({
	highlights,
}: {
	highlights?: Array<CredibilityHighlight>;
}) {
	if (!highlights || highlights.length === 0) {
		return null;
	}

	return (
		<div className="mt-4">
			<h4 className="text-sm font-semibold text-gray-700 mb-2">
				Employer Feedback
			</h4>
			<div className="space-y-3">
				{highlights.map((highlight) => (
					<div
						key={highlight.quote || Math.random()}
						className={`p-3 rounded-lg border-l-4 ${
							highlight.sentiment === "POSITIVE"
								? "bg-green-50 border-green-400"
								: "bg-gray-50 border-gray-400"
						}`}
					>
						<blockquote className="text-sm text-gray-700 italic">
							"{highlight.quote}"
						</blockquote>
					</div>
				))}
			</div>
		</div>
	);
}

