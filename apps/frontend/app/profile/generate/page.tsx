"use client";

import { useState, useCallback } from "react";
import { ReferenceLetterUpload } from "@/components/profile/ReferenceLetterUpload";
import { GenerateProfileButton } from "@/components/profile/GenerateProfileButton";
import { ProfileEditor } from "@/components/profile/ProfileEditor";
import type { ReferenceLetter } from "@/lib/api/referenceLetters";
import type { Profile } from "@/lib/api/profile";

export default function GenerateProfilePage() {
	const [uploadedLetters, setUploadedLetters] = useState<ReferenceLetter[]>([]);
	const [generatedProfile, setGeneratedProfile] = useState<Profile | null>(
		null,
	);
	const [error, setError] = useState<string | null>(null);
	const [step, setStep] = useState<"upload" | "generate" | "edit">("upload");

	const handleUploadComplete = useCallback((letter: ReferenceLetter) => {
		setUploadedLetters((prev) => [...prev, letter]);
		setError(null);
	}, []);

	const handleGenerateComplete = useCallback((profile: Profile) => {
		setGeneratedProfile(profile);
		setStep("edit");
		setError(null);
	}, []);

	const handleError = useCallback((errorMsg: string) => {
		setError(errorMsg);
	}, []);

	return (
		<div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
			<div className="max-w-4xl mx-auto px-4 py-12">
				{/* Header */}
				<div className="text-center mb-12">
					<h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
						Generate Your Smart Profile
					</h1>
					<p className="text-lg text-gray-600 dark:text-gray-400">
						Upload your reference letters and let AI create a professional
						profile with credibility highlights
					</p>
				</div>

				{/* Progress Steps */}
				<div className="flex justify-center mb-12">
					<div className="flex items-center gap-4">
						<StepIndicator
							number={1}
							label="Upload"
							isActive={step === "upload"}
							isComplete={uploadedLetters.length > 0}
						/>
						<div className="w-12 h-0.5 bg-gray-300 dark:bg-gray-600" />
						<StepIndicator
							number={2}
							label="Generate"
							isActive={step === "generate"}
							isComplete={generatedProfile !== null}
						/>
						<div className="w-12 h-0.5 bg-gray-300 dark:bg-gray-600" />
						<StepIndicator
							number={3}
							label="Edit"
							isActive={step === "edit"}
							isComplete={false}
						/>
					</div>
				</div>

				{/* Error Display */}
				{error && (
					<div className="mb-8 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl">
						<p className="text-red-600 dark:text-red-400">{error}</p>
					</div>
				)}

				{/* Content based on step */}
				{step === "upload" && (
					<div className="space-y-8">
						<div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
							<h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-6">
								Upload Reference Letters
							</h2>
							<ReferenceLetterUpload
								onUploadComplete={handleUploadComplete}
								onError={handleError}
							/>
						</div>

						{uploadedLetters.length > 0 && (
							<div className="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700">
								<h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
									Uploaded Letters ({uploadedLetters.length})
								</h2>
								<ul className="space-y-2 mb-6">
									{uploadedLetters.map((letter) => (
										<li
											key={letter.id}
											className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg"
										>
											<div className="w-8 h-8 rounded bg-green-100 dark:bg-green-900 flex items-center justify-center">
												<svg
													className="w-4 h-4 text-green-600 dark:text-green-400"
													fill="none"
													stroke="currentColor"
													viewBox="0 0 24 24"
													aria-hidden="true"
												>
													<title>Checkmark</title>
													<path
														strokeLinecap="round"
														strokeLinejoin="round"
														strokeWidth={2}
														d="M5 13l4 4L19 7"
													/>
												</svg>
											</div>
											<span className="text-gray-700 dark:text-gray-300">
												{letter.fileName}
											</span>
											<span className="text-sm text-gray-500 ml-auto">
												{letter.status}
											</span>
										</li>
									))}
								</ul>

								<div className="flex justify-center">
									<GenerateProfileButton
										onGenerateComplete={handleGenerateComplete}
										onError={handleError}
										disabled={uploadedLetters.length === 0}
									/>
								</div>
							</div>
						)}
					</div>
				)}

				{step === "edit" && generatedProfile && (
					<div className="space-y-8">
						<div className="flex items-center justify-between">
							<h2 className="text-2xl font-semibold text-gray-900 dark:text-white">
								Edit Your Profile
							</h2>
							<button
								type="button"
								onClick={() => setStep("upload")}
								className="text-blue-600 hover:text-blue-700 font-medium"
							>
								← Upload More Letters
							</button>
						</div>

						<ProfileEditor
							profile={generatedProfile}
							onUpdate={setGeneratedProfile}
							onError={handleError}
						/>
					</div>
				)}
			</div>
		</div>
	);
}

interface StepIndicatorProps {
	number: number;
	label: string;
	isActive: boolean;
	isComplete: boolean;
}

function StepIndicator({
	number,
	label,
	isActive,
	isComplete,
}: StepIndicatorProps) {
	return (
		<div className="flex flex-col items-center gap-2">
			<div
				className={`w-10 h-10 rounded-full flex items-center justify-center font-semibold transition-colors ${
					isComplete
						? "bg-green-500 text-white"
						: isActive
							? "bg-blue-600 text-white"
							: "bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400"
				}`}
			>
				{isComplete ? (
					<svg
						className="w-5 h-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						aria-hidden="true"
					>
						<title>Complete</title>
						<path
							strokeLinecap="round"
							strokeLinejoin="round"
							strokeWidth={2}
							d="M5 13l4 4L19 7"
						/>
					</svg>
				) : (
					number
				)}
			</div>
			<span
				className={`text-sm font-medium ${
					isActive || isComplete
						? "text-gray-900 dark:text-white"
						: "text-gray-500 dark:text-gray-400"
				}`}
			>
				{label}
			</span>
		</div>
	);
}
