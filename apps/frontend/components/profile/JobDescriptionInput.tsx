"use client";

const MIN_LENGTH = 50;
const MAX_LENGTH = 10000;

interface JobDescriptionInputProps {
	value: string;
	onChange: (value: string) => void;
	onTailor?: () => void;
	isLoading?: boolean;
	error?: string | null;
	disabled?: boolean;
}

export const JobDescriptionInput = ({
	value,
	onChange,
	onTailor,
	isLoading = false,
	error = null,
	disabled = false,
}: JobDescriptionInputProps) => {
	const trimmedLength = value.trim().length;
	const isTooShort = trimmedLength > 0 && trimmedLength < MIN_LENGTH;
	const isTooLong = trimmedLength > MAX_LENGTH;
	const isValid = trimmedLength >= MIN_LENGTH && trimmedLength <= MAX_LENGTH;

	const getCharacterCountColor = () => {
		if (isTooLong) return "text-red-500";
		if (isTooShort) return "text-amber-500";
		if (isValid) return "text-green-600";
		return "text-gray-500";
	};

	const getCharacterCountText = () => {
		if (trimmedLength === 0) {
			return "Copy and paste the complete job description for best results";
		}
		if (isTooShort) {
			return `${trimmedLength}/${MIN_LENGTH} characters (minimum ${MIN_LENGTH} required)`;
		}
		if (isTooLong) {
			return `${trimmedLength}/${MAX_LENGTH} characters (maximum exceeded)`;
		}
		return `${trimmedLength} characters ✓`;
	};

	return (
		<div className="space-y-4">
			<div className="space-y-2">
				<label
					htmlFor="job-description"
					className="block text-sm font-medium text-gray-700"
				>
					Job Description
				</label>
				<textarea
					id="job-description"
					value={value}
					onChange={(e) => onChange(e.target.value)}
					disabled={disabled || isLoading}
					placeholder="Paste the complete job description here...

Example:
We are looking for a Senior Software Engineer with experience in Go, microservices architecture, and cloud platforms. The ideal candidate has 5+ years of experience building scalable web applications..."
					className={`w-full min-h-[200px] p-4 border rounded-xl resize-y
						${disabled || isLoading ? "bg-gray-100 cursor-not-allowed text-gray-500" : "bg-white"}
						${error ? "border-red-300 focus:border-red-500 focus:ring-red-200" : "border-gray-200 focus:border-indigo-500 focus:ring-indigo-200"}
						focus:ring-2 transition-colors outline-none text-gray-700 placeholder-gray-400`}
				/>
				<p className={`text-xs ${getCharacterCountColor()}`}>
					{getCharacterCountText()}
				</p>
			</div>

			{error && (
				<div className="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
					<svg
						className="w-5 h-5 flex-shrink-0"
						fill="currentColor"
						viewBox="0 0 20 20"
					>
						<title>Error</title>
						<path
							fillRule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
							clipRule="evenodd"
						/>
					</svg>
					<span>{error}</span>
				</div>
			)}

			{onTailor && (
				<button
					type="button"
					onClick={onTailor}
					disabled={!isValid || isLoading || disabled}
					className={`w-full py-3 px-6 rounded-xl font-semibold text-white
						transition-all duration-300 transform
						${
							!isValid || isLoading || disabled
								? "bg-gray-400 cursor-not-allowed"
								: "bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 hover:scale-[1.02] hover:shadow-lg"
						}`}
				>
					{isLoading ? (
						<span className="flex items-center justify-center gap-3">
							<svg
								className="animate-spin h-5 w-5"
								fill="none"
								viewBox="0 0 24 24"
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
							<span>Analyzing job description with AI...</span>
						</span>
					) : (
						<span className="flex items-center justify-center gap-2">
							<svg
								className="w-5 h-5"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<title>Tailor</title>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
								/>
							</svg>
							Tailor My Profile
						</span>
					)}
				</button>
			)}

			{isLoading && (
				<div className="text-center text-sm text-gray-500">
					<p>
						This may take 10-30 seconds while our AI analyzes your profile
						against the job requirements...
					</p>
				</div>
			)}
		</div>
	);
};
