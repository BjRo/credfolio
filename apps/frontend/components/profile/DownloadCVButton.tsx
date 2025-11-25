"use client";

import { useState } from "react";
import { downloadCV } from "@/lib/api/profile";

interface DownloadCVButtonProps {
	profileId: string;
	tailoredId?: string;
	disabled?: boolean;
	onError?: (error: string) => void;
}

export const DownloadCVButton = ({
	profileId,
	tailoredId,
	disabled = false,
	onError,
}: DownloadCVButtonProps) => {
	const [isDownloading, setIsDownloading] = useState(false);

	const handleDownload = async () => {
		setIsDownloading(true);

		try {
			const blob = await downloadCV(profileId, tailoredId);

			// Create download link
			const url = URL.createObjectURL(blob);
			const a = document.createElement("a");
			a.href = url;
			a.download = tailoredId ? "tailored-cv.pdf" : "cv.pdf";
			document.body.appendChild(a);
			a.click();
			document.body.removeChild(a);
			URL.revokeObjectURL(url);
		} catch (err) {
			const errorMessage = "Failed to download CV";
			if (onError) {
				onError(errorMessage);
			}
			console.error("Download error:", err);
		} finally {
			setIsDownloading(false);
		}
	};

	return (
		<button
			type="button"
			onClick={handleDownload}
			disabled={disabled || isDownloading}
			className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all ${
				disabled || isDownloading
					? "bg-gray-200 text-gray-400 cursor-not-allowed"
					: "bg-gradient-to-r from-emerald-500 to-teal-500 text-white hover:from-emerald-600 hover:to-teal-600 shadow-md hover:shadow-lg"
			}`}
		>
			{isDownloading ? (
				<>
					<svg
						className="animate-spin w-5 h-5"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<circle
							className="opacity-25"
							cx="12"
							cy="12"
							r="10"
							stroke="currentColor"
							strokeWidth="4"
							fill="none"
						/>
						<path
							className="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
						/>
					</svg>
					Downloading...
				</>
			) : (
				<>
					<svg
						className="w-5 h-5"
						fill="currentColor"
						viewBox="0 0 20 20"
						aria-hidden="true"
					>
						<path
							fillRule="evenodd"
							d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
							clipRule="evenodd"
						/>
					</svg>
					Download CV
				</>
			)}
		</button>
	);
};
