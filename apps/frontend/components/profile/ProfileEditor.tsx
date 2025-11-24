"use client";

import { useState } from "react";
import { DefaultService } from "../../lib/api/generated";
import type { Profile } from "../../lib/api/generated/models/Profile";
import { getErrorMessage } from "../../lib/utils/errorMessages";

const MAX_LENGTH = 2000;

export default function ProfileEditor({
	profile,
	onUpdate,
}: {
	profile: Profile;
	onUpdate: (profile: Profile) => void;
}) {
	const [summary, setSummary] = useState(profile.summary || "");
	const [updating, setUpdating] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const getCharacterCount = (text: string) => {
		// Count characters using proper Unicode support
		return Array.from(text).length;
	};

	const validateInput = (text: string): string | null => {
		const charCount = getCharacterCount(text);
		if (charCount > MAX_LENGTH) {
			return `Summary must be at most ${MAX_LENGTH} characters (currently ${charCount})`;
		}
		return null;
	};

	const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		const value = e.target.value;
		setSummary(value);
		// Clear error when user starts typing
		if (error) {
			const validationError = validateInput(value);
			setError(validationError);
		}
	};

	const handleSave = async () => {
		const validationError = validateInput(summary);
		if (validationError) {
			setError(validationError);
			return;
		}

		setUpdating(true);
		setError(null);
		try {
			const updated = await DefaultService.updateProfile({
				summary,
			});
			onUpdate(updated);
			alert("Profile updated!");
		} catch (err) {
			console.error(err);
			alert(getErrorMessage(err));
		} finally {
			setUpdating(false);
		}
	};

	const charCount = getCharacterCount(summary);
	const isInvalid = charCount > MAX_LENGTH;
	const showWarning = charCount > MAX_LENGTH * 0.9; // Show warning at 90% of max

	return (
		<div className="bg-white p-4 rounded shadow mt-4">
			<h3 className="text-lg font-bold mb-2">Edit Profile</h3>
			<div className="mb-4">
				<label
					htmlFor="summary-input"
					className="block text-sm font-medium text-gray-700"
				>
					Summary
				</label>
				<textarea
					id="summary-input"
					className={`mt-1 block w-full border rounded-md shadow-sm p-2 ${
						error || isInvalid
							? "border-red-300 focus:border-red-500 focus:ring-red-500"
							: "border-gray-300"
					}`}
					rows={4}
					value={summary}
					onChange={handleChange}
					disabled={updating}
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
					</div>
				</div>
			</div>
			<button
				type="button"
				onClick={handleSave}
				disabled={updating || isInvalid}
				className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
			>
				{updating ? "Saving..." : "Save Changes"}
			</button>
		</div>
	);
}
