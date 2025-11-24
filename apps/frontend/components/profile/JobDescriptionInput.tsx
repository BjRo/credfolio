"use client";

import { useState } from "react";

interface JobDescriptionInputProps {
	onSubmit: (jobDescription: string) => void;
	isLoading?: boolean;
}

export default function JobDescriptionInput({
	onSubmit,
	isLoading = false,
}: JobDescriptionInputProps) {
	const [jobDescription, setJobDescription] = useState("");

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (jobDescription.trim()) {
			onSubmit(jobDescription.trim());
		}
	};

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
				onChange={(e) => setJobDescription(e.target.value)}
				placeholder="Paste the job description here..."
				rows={8}
				className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none"
				disabled={isLoading}
			/>
			<button
				type="submit"
				disabled={!jobDescription.trim() || isLoading}
				className="mt-3 px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
			>
				{isLoading ? "Tailoring..." : "Tailor Profile"}
			</button>
		</form>
	);
}
