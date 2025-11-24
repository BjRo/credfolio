"use client";

import { useState } from "react";

interface JobDescriptionInputProps {
	onSubmit: (jobDescription: string) => void;
	isLoading?: boolean;
}

const MIN_LENGTH = 10;
const MAX_LENGTH = 10000;

export default function JobDescriptionInput({
	onSubmit,
	isLoading = false,
}: JobDescriptionInputProps) {
	const [jobDescription, setJobDescription] = useState("");
	const [error, setError] = useState<string | null>(null);

	const getCharacterCount = (text: string) => {
		// Count characters using proper Unicode support
		return Array.from(text).length;
	};

	const validateInput = (text: string): string | null => {
		const trimmed = text.trim();
		const charCount = getCharacterCount(trimmed);

		if (trimmed.length === 0) {
			return "Job description is required";
		}
		if (charCount < MIN_LENGTH) {
			return `Job description must be at least ${MIN_LENGTH} characters (currently ${charCount})`;
		}
		if (charCount > MAX_LENGTH) {
			return `Job description must be at most ${MAX_LENGTH} characters (currently ${charCount})`;
		}
		return null;
	};

	const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		const value = e.target.value;
		setJobDescription(value);
		// Clear error when user starts typing
		if (error) {
			const validationError = validateInput(value);
			setError(validationError);
		}
	};

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		const trimmed = jobDescription.trim();
		const validationError = validateInput(trimmed);

		if (validationError) {
			setError(validationError);
			return;
		}

		setError(null);
		onSubmit(trimmed);
	};

	const charCount = getCharacterCount(jobDescription);
	const isInvalid = charCount < MIN_LENGTH || charCount > MAX_LENGTH;
	const showWarning = charCount > MAX_LENGTH * 0.9; // Show warning at 90% of max

	return (
		<form onSubmit={handleSubmit} className="mb-6">
			<label
				htmlFor="job-description"
				className="block text-sm font-medium text-gray-700 mb-2"
			>
				Job Description
			</label>
			<textarea
				id="job-description"
				value={jobDescription}
				onChange={handleChange}
				placeholder="Paste the job description here..."
				rows={8}
				className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none ${
					error || isInvalid
						? "border-red-300 focus:border-red-500 focus:ring-red-500"
						: "border-gray-300"
				}`}
				disabled={isLoading}
			/>
			<div className="mt-2 flex justify-between items-center">
				<div className="flex-1">
					{error && <p className="text-sm text-red-600 mt-1">{error}</p>}
				</div>
				<div
					className={`text-sm ${
						showWarning || isInvalid ? "text-red-600" : "text-gray-500"
					}`}
				>
					{charCount} / {MAX_LENGTH} characters
					{charCount < MIN_LENGTH && (
						<span className="ml-1">(min {MIN_LENGTH})</span>
					)}
				</div>
			</div>
			<button
				type="submit"
				disabled={!jobDescription.trim() || isLoading || !!error || isInvalid}
				className="mt-3 px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
			>
				{isLoading ? "Tailoring..." : "Tailor Profile"}
			</button>
		</form>
	);
}
