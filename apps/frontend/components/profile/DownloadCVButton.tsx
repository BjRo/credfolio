"use client";

import { useState } from "react";
import { downloadCV } from "../../lib/api/profile";
import { getErrorMessage } from "../../lib/utils/errorMessages";

interface DownloadCVButtonProps {
	profileId: string;
	jobMatchId?: string;
	className?: string;
}

export default function DownloadCVButton({
	profileId,
	jobMatchId,
	className = "",
}: DownloadCVButtonProps) {
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleDownload = async () => {
		try {
			setLoading(true);
			setError(null);

			const blob = await downloadCV(profileId, jobMatchId);

			// Create download link
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement("a");
			a.href = url;
			a.download = jobMatchId ? "tailored-cv.pdf" : "cv.pdf";
			document.body.appendChild(a);
			a.click();

			// Cleanup
			window.URL.revokeObjectURL(url);
			document.body.removeChild(a);
		} catch (err) {
			console.error("Failed to download CV:", err);
			setError(getErrorMessage(err));
		} finally {
			setLoading(false);
		}
	};

	return (
		<div>
			<button
				type="button"
				onClick={handleDownload}
				disabled={loading}
				className={`px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center gap-2 ${className}`}
			>
				{loading ? (
					<>
						<div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
						<span>Downloading...</span>
					</>
				) : (
					<>
						<svg
							className="w-5 h-5"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							role="img"
							aria-label="Download icon"
						>
							<title>Download icon</title>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
						<span>Download CV</span>
					</>
				)}
			</button>
			{error && <p className="mt-2 text-sm text-red-600">{error}</p>}
		</div>
	);
}
