"use client";

import { useState } from "react";
import { generateProfile, type Profile } from "@/lib/api/profile";

interface GenerateProfileButtonProps {
	onGenerateComplete?: (profile: Profile) => void;
	onError?: (error: string) => void;
	disabled?: boolean;
}

export function GenerateProfileButton({
	onGenerateComplete,
	onError,
	disabled = false,
}: GenerateProfileButtonProps) {
	const [isGenerating, setIsGenerating] = useState(false);

	const handleGenerate = async () => {
		setIsGenerating(true);

		try {
			const profile = await generateProfile();
			onGenerateComplete?.(profile);
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : "Generation failed";
			onError?.(errorMsg);
		} finally {
			setIsGenerating(false);
		}
	};

	return (
		<button
			type="button"
			onClick={handleGenerate}
			disabled={disabled || isGenerating}
			className={`
        relative px-6 py-3 rounded-xl font-semibold text-white
        transition-all duration-300 transform
        ${
					disabled || isGenerating
						? "bg-gray-400 cursor-not-allowed"
						: "bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 hover:scale-105 hover:shadow-lg"
				}
      `}
		>
			{isGenerating ? (
				<span className="flex items-center gap-2">
					<svg
						className="animate-spin h-5 w-5"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<title>Loading</title>
						<circle
							className="opacity-25"
							cx="12"
							cy="12"
							r="10"
							stroke="currentColor"
							strokeWidth="4"
						/>
						<path
							className="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						/>
					</svg>
					Generating Profile...
				</span>
			) : (
				<span className="flex items-center gap-2">
					<svg
						className="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<title>Generate</title>
						<path
							strokeLinecap="round"
							strokeLinejoin="round"
							strokeWidth={2}
							d="M13 10V3L4 14h7v7l9-11h-7z"
						/>
					</svg>
					Generate Profile with AI
				</span>
			)}
		</button>
	);
}
