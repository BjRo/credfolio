"use client";

import LoadingSpinner from "./LoadingSpinner";

interface ProgressIndicatorProps {
	message: string;
	subMessage?: string;
	showSpinner?: boolean;
	className?: string;
}

export default function ProgressIndicator({
	message,
	subMessage,
	showSpinner = true,
	className = "",
}: ProgressIndicatorProps) {
	return (
		<div
			className={`flex flex-col items-center justify-center py-8 ${className}`}
			aria-live="polite"
			aria-busy="true"
		>
			{showSpinner && (
				<div className="mb-4">
					<LoadingSpinner size="lg" className="text-indigo-600" />
				</div>
			)}
			<p className="text-lg font-medium text-gray-900">{message}</p>
			{subMessage && <p className="mt-2 text-sm text-gray-600">{subMessage}</p>}
		</div>
	);
}
