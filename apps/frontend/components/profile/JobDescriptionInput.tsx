"use client";

interface JobDescriptionInputProps {
	value: string;
	onChange: (value: string) => void;
	disabled?: boolean;
}

export const JobDescriptionInput = ({
	value,
	onChange,
	disabled = false,
}: JobDescriptionInputProps) => {
	return (
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
				disabled={disabled}
				placeholder="Paste the job description here to tailor your profile..."
				className={`w-full min-h-[200px] p-4 border rounded-xl resize-y
					${disabled ? "bg-gray-100 cursor-not-allowed text-gray-500" : "bg-white"}
					border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200
					transition-colors outline-none text-gray-700 placeholder-gray-400`}
			/>
			<p className="text-xs text-gray-500">
				{value.length > 0
					? `${value.length} characters`
					: "Copy and paste the complete job description for best results"}
			</p>
		</div>
	);
};
