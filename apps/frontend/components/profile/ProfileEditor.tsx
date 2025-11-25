"use client";

import { useState, useCallback } from "react";
import {
	updateProfile,
	type Profile,
	type WorkExperience,
} from "@/lib/api/profile";

interface ProfileEditorProps {
	profile: Profile;
	onUpdate?: (profile: Profile) => void;
	onError?: (error: string) => void;
}

export function ProfileEditor({
	profile,
	onUpdate,
	onError,
}: ProfileEditorProps) {
	const [editedProfile, setEditedProfile] = useState(profile);
	const [isSaving, setIsSaving] = useState(false);
	const [editingField, setEditingField] = useState<string | null>(null);

	const handleSummaryChange = useCallback((value: string) => {
		setEditedProfile((prev) => ({ ...prev, summary: value }));
	}, []);

	const handleSave = useCallback(async () => {
		setIsSaving(true);
		try {
			const updated = await updateProfile({ summary: editedProfile.summary });
			onUpdate?.(updated);
			setEditingField(null);
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : "Save failed";
			onError?.(errorMsg);
		} finally {
			setIsSaving(false);
		}
	}, [editedProfile.summary, onUpdate, onError]);

	return (
		<div className="space-y-6">
			{/* Summary Section */}
			<div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700">
				<div className="flex items-center justify-between mb-4">
					<h2 className="text-xl font-semibold text-gray-900 dark:text-white">
						Professional Summary
					</h2>
					{editingField !== "summary" && (
						<button
							type="button"
							onClick={() => setEditingField("summary")}
							className="text-blue-600 hover:text-blue-700 text-sm font-medium"
						>
							Edit
						</button>
					)}
				</div>

				{editingField === "summary" ? (
					<div className="space-y-3">
						<textarea
							value={editedProfile.summary}
							onChange={(e) => handleSummaryChange(e.target.value)}
							className="w-full p-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-white resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							rows={4}
							placeholder="Write a professional summary..."
						/>
						<div className="flex gap-2">
							<button
								type="button"
								onClick={handleSave}
								disabled={isSaving}
								className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
							>
								{isSaving ? "Saving..." : "Save"}
							</button>
							<button
								type="button"
								onClick={() => {
									setEditedProfile(profile);
									setEditingField(null);
								}}
								className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
							>
								Cancel
							</button>
						</div>
					</div>
				) : (
					<p className="text-gray-600 dark:text-gray-300 whitespace-pre-wrap">
						{editedProfile.summary || "No summary yet. Click Edit to add one."}
					</p>
				)}
			</div>

			{/* Work Experience Section */}
			<div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700">
				<h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
					Work Experience
				</h2>

				{editedProfile.workExperiences.length > 0 ? (
					<div className="space-y-4">
						{editedProfile.workExperiences.map((exp) => (
							<WorkExperienceCard key={exp.id} experience={exp} />
						))}
					</div>
				) : (
					<p className="text-gray-500 dark:text-gray-400">
						No work experience added yet. Upload reference letters to extract
						experience.
					</p>
				)}
			</div>

			{/* Skills Section */}
			<div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700">
				<h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
					Skills
				</h2>

				{editedProfile.skills.length > 0 ? (
					<div className="flex flex-wrap gap-2">
						{editedProfile.skills.map((skill) => (
							<span
								key={skill}
								className="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded-full text-sm font-medium"
							>
								{skill}
							</span>
						))}
					</div>
				) : (
					<p className="text-gray-500 dark:text-gray-400">
						No skills added yet. Skills will be extracted from reference
						letters.
					</p>
				)}
			</div>
		</div>
	);
}

interface WorkExperienceCardProps {
	experience: WorkExperience;
}

function WorkExperienceCard({ experience }: WorkExperienceCardProps) {
	const [isExpanded, setIsExpanded] = useState(false);

	return (
		<div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
			<div className="flex justify-between items-start">
				<div>
					<h3 className="font-semibold text-gray-900 dark:text-white">
						{experience.role}
					</h3>
					<p className="text-gray-600 dark:text-gray-400">
						{experience.companyName}
					</p>
					<p className="text-sm text-gray-500 dark:text-gray-500">
						{experience.startDate} - {experience.endDate || "Present"}
					</p>
				</div>
				<button
					type="button"
					onClick={() => setIsExpanded(!isExpanded)}
					className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
					aria-label={isExpanded ? "Collapse details" : "Expand details"}
				>
					<svg
						className={`w-5 h-5 transition-transform ${isExpanded ? "rotate-180" : ""}`}
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<title>Toggle details</title>
						<path
							strokeLinecap="round"
							strokeLinejoin="round"
							strokeWidth={2}
							d="M19 9l-7 7-7-7"
						/>
					</svg>
				</button>
			</div>

			{isExpanded && (
				<div className="mt-4 space-y-3">
					{experience.description && (
						<p className="text-gray-600 dark:text-gray-300 text-sm">
							{experience.description}
						</p>
					)}

					{experience.credibilityHighlights &&
						experience.credibilityHighlights.length > 0 && (
							<div className="mt-3">
								<h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
									Credibility Highlights
								</h4>
								<div className="space-y-2">
									{experience.credibilityHighlights.map((highlight) => (
										<div
											key={highlight.quote}
											className="bg-green-50 dark:bg-green-900/20 border-l-4 border-green-500 p-3 rounded-r"
										>
											<p className="text-gray-700 dark:text-gray-300 text-sm italic">
												"{highlight.quote}"
											</p>
											<span className="text-xs text-green-600 dark:text-green-400 font-medium mt-1 inline-block">
												{highlight.sentiment}
											</span>
										</div>
									))}
								</div>
							</div>
						)}
				</div>
			)}
		</div>
	);
}
