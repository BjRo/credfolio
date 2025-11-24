"use client";

import { useState } from "react";
import { generateProfile } from "../../lib/api/profile";
import type { Profile } from "../../lib/api/generated/models/Profile";
import { getErrorMessage } from "../../lib/utils/errorMessages";
import LoadingSpinner from "./LoadingSpinner";

export default function GenerateProfileButton({
	onGenerateComplete,
}: {
	onGenerateComplete?: (profile: Profile) => void;
}) {
	const [generating, setGenerating] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleGenerate = async () => {
		setGenerating(true);
		setError(null);

		try {
			const profile = await generateProfile();
			if (onGenerateComplete) onGenerateComplete(profile);
		} catch (err) {
			setError(getErrorMessage(err));
			console.error(err);
		} finally {
			setGenerating(false);
		}
	};

	return (
		<div className="mt-4">
			<button
				type="button"
				onClick={handleGenerate}
				disabled={generating}
				className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-blue-300 disabled:cursor-not-allowed flex items-center gap-2 transition-colors"
			>
				{generating && <LoadingSpinner size="sm" className="text-white" />}
				<span>
					{generating ? "Generating Profile..." : "Generate Smart Profile"}
				</span>
			</button>
			{generating && (
				<p className="mt-2 text-sm text-gray-600">
					This may take up to 60 seconds. Please wait...
				</p>
			)}
			{error && <p className="mt-2 text-sm text-red-600">{error}</p>}
		</div>
	);
}
