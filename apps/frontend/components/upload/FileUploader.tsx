"use client";

import { useState } from "react";
import { uploadReferenceLetter } from "../../services/api";

export default function FileUploader() {
	const [status, setStatus] = useState<
		"idle" | "uploading" | "success" | "error"
	>("idle");
	const [message, setMessage] = useState("");

	const handleFile = async (file: File) => {
		// Allow PDF, plain text, and Markdown
		const validTypes = ["application/pdf", "text/plain", "text/markdown"];
		const validExtensions = [".pdf", ".txt", ".md"];

		const hasValidType = validTypes.includes(file.type);
		const hasValidExt = validExtensions.some((ext) =>
			file.name.toLowerCase().endsWith(ext),
		);

		if (!hasValidType && !hasValidExt) {
			setStatus("error");
			setMessage("Only PDF, Text, or Markdown files are supported");
			return;
		}

		setStatus("uploading");
		try {
			const res = await uploadReferenceLetter(file);
			setStatus("success");
			setMessage("File uploaded successfully! Processing...");
			// Optional: redirect to dashboard or show link
			// window.location.href = `/dashboard?user_id=${res.user_id}`;
		} catch (err) {
			setStatus("error");
			setMessage("Upload failed. Please try again.");
		}
	};

	const onDrop = (e: React.DragEvent) => {
		e.preventDefault();
		if (e.dataTransfer.files?.[0]) {
			handleFile(e.dataTransfer.files[0]);
		}
	};

	const onDragOver = (e: React.DragEvent) => {
		e.preventDefault();
	};

	const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		if (e.target.files?.[0]) {
			handleFile(e.target.files[0]);
		}
	};

	return (
		<div className="w-full max-w-md mx-auto">
			<div
				onDrop={onDrop}
				onDragOver={onDragOver}
				className="border-2 border-dashed border-gray-300 rounded-lg p-10 text-center cursor-pointer hover:border-blue-500 transition-colors bg-white dark:bg-gray-800"
			>
				<input
					type="file"
					accept=".pdf,.txt,.md"
					className="hidden"
					id="file-upload"
					onChange={onChange}
				/>
				<label
					htmlFor="file-upload"
					className="cursor-pointer block w-full h-full"
				>
					<div className="text-lg mb-2 font-medium">
						Drag & Drop or Click to Upload
					</div>
					<div className="text-sm text-gray-500">PDF, Text, or Markdown</div>
				</label>
			</div>

			{status === "uploading" && (
				<div className="mt-4 text-blue-600 text-center animate-pulse">
					Uploading...
				</div>
			)}
			{status === "success" && (
				<div className="mt-4 text-green-600 text-center">{message}</div>
			)}
			{status === "error" && (
				<div className="mt-4 text-red-600 text-center">{message}</div>
			)}
		</div>
	);
}
