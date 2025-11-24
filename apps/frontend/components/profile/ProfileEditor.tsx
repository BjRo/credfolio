"use client";

import { useState } from "react";
import { DefaultService } from "../../lib/api/generated";
import type { Profile } from "../../lib/api/generated/models/Profile";
import { getErrorMessage } from "../../lib/utils/errorMessages";

export default function ProfileEditor({
	profile,
	onUpdate,
}: {
	profile: Profile;
	onUpdate: (profile: Profile) => void;
}) {
	const [summary, setSummary] = useState(profile.summary || "");
	const [updating, setUpdating] = useState(false);

	const handleSave = async () => {
		setUpdating(true);
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
					className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
					rows={4}
					value={summary}
					onChange={(e) => setSummary(e.target.value)}
				/>
			</div>
			<button
				type="button"
				onClick={handleSave}
				disabled={updating}
				className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700"
			>
				{updating ? "Saving..." : "Save Changes"}
			</button>
		</div>
	);
}
