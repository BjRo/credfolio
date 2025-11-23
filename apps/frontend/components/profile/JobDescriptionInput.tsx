"use client";

import { useState } from "react";
import type { TailoringResult } from "../../types";
import { tailorProfile } from "../../services/api";

export default function JobDescriptionInput({
	userID,
	onTailored,
}: { userID: string; onTailored: (result: TailoringResult) => void }) {
	const [jd, setJd] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState("");

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (!jd.trim()) return;

		setLoading(true);
		setError("");
		try {
			const res = await tailorProfile(userID, jd);
			onTailored(res);
		} catch (err) {
			setError("Failed to tailor profile");
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="mb-8 p-6 bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
			<h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
				Tailor CV to Job
			</h3>
			<form onSubmit={handleSubmit}>
				<textarea
					className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
					rows={4}
					placeholder="Paste job description here..."
					value={jd}
					onChange={(e) => setJd(e.target.value)}
				/>
				{error && <p className="text-red-500 mt-2 text-sm">{error}</p>}
				<button
					type="submit"
					disabled={loading || !jd.trim()}
					className="mt-3 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
				>
					{loading ? "Analyzing..." : "Tailor Profile"}
				</button>
			</form>
		</div>
	);
}
