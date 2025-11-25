"use client";

import { useState, useCallback } from "react";
import {
	uploadReferenceLetter,
	type ReferenceLetter,
} from "@/lib/api/referenceLetters";

interface ReferenceLetterUploadProps {
	onUploadComplete?: (letter: ReferenceLetter) => void;
	onError?: (error: string) => void;
}

export function ReferenceLetterUpload({
	onUploadComplete,
	onError,
}: ReferenceLetterUploadProps) {
	const [selectedFile, setSelectedFile] = useState<File | null>(null);
	const [isUploading, setIsUploading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [isDragOver, setIsDragOver] = useState(false);

	const handleFileSelect = useCallback(
		(file: File) => {
			const allowedExtensions = [".txt", ".md", ".markdown"];
			const ext = file.name.toLowerCase().substring(file.name.lastIndexOf("."));

			if (!allowedExtensions.includes(ext)) {
				const errorMsg = "Please select a .txt or .md file";
				setError(errorMsg);
				onError?.(errorMsg);
				return;
			}

			setSelectedFile(file);
			setError(null);
		},
		[onError],
	);

	const handleDrop = useCallback(
		(e: React.DragEvent) => {
			e.preventDefault();
			setIsDragOver(false);

			const file = e.dataTransfer.files[0];
			if (file) {
				handleFileSelect(file);
			}
		},
		[handleFileSelect],
	);

	const handleDragOver = useCallback((e: React.DragEvent) => {
		e.preventDefault();
		setIsDragOver(true);
	}, []);

	const handleDragLeave = useCallback(() => {
		setIsDragOver(false);
	}, []);

	const handleInputChange = useCallback(
		(e: React.ChangeEvent<HTMLInputElement>) => {
			const file = e.target.files?.[0];
			if (file) {
				handleFileSelect(file);
			}
		},
		[handleFileSelect],
	);

	const handleUpload = useCallback(async () => {
		if (!selectedFile) return;

		setIsUploading(true);
		setError(null);

		try {
			const letter = await uploadReferenceLetter(selectedFile);
			setSelectedFile(null);
			onUploadComplete?.(letter);
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : "Upload failed";
			setError(errorMsg);
			onError?.(errorMsg);
		} finally {
			setIsUploading(false);
		}
	}, [selectedFile, onUploadComplete, onError]);

	return (
		<div className="w-full">
			<div
				className={`border-2 border-dashed rounded-xl p-8 text-center transition-all duration-200 ${
					isDragOver
						? "border-blue-500 bg-blue-50 dark:bg-blue-950/20"
						: "border-gray-300 dark:border-gray-700 hover:border-gray-400 dark:hover:border-gray-600"
				}`}
				onDrop={handleDrop}
				onDragOver={handleDragOver}
				onDragLeave={handleDragLeave}
			>
				<input
					type="file"
					accept=".txt,.md,.markdown"
					onChange={handleInputChange}
					className="hidden"
					id="file-upload"
					disabled={isUploading}
				/>

				<label
					htmlFor="file-upload"
					className="cursor-pointer flex flex-col items-center gap-3"
				>
					<div className="w-12 h-12 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
						<svg
							className="w-6 h-6 text-gray-500"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<title>Upload</title>
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth={2}
								d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
							/>
						</svg>
					</div>

					<div>
						<span className="text-blue-600 dark:text-blue-400 font-medium">
							Click to upload
						</span>
						<span className="text-gray-500 dark:text-gray-400">
							{" "}
							or drag and drop
						</span>
					</div>

					<p className="text-sm text-gray-500 dark:text-gray-400">
						.txt or .md files only
					</p>
				</label>
			</div>

			{selectedFile && (
				<div className="mt-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg flex items-center justify-between">
					<div className="flex items-center gap-3">
						<div className="w-10 h-10 rounded bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
							<svg
								className="w-5 h-5 text-blue-600 dark:text-blue-400"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
								aria-hidden="true"
							>
								<title>Document</title>
								<path
									strokeLinecap="round"
									strokeLinejoin="round"
									strokeWidth={2}
									d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
								/>
							</svg>
						</div>
						<div>
							<p className="font-medium text-gray-900 dark:text-white">
								{selectedFile.name}
							</p>
							<p className="text-sm text-gray-500">
								{(selectedFile.size / 1024).toFixed(1)} KB
							</p>
						</div>
					</div>

					<button
						type="button"
						onClick={handleUpload}
						disabled={isUploading}
						className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
					>
						{isUploading ? "Uploading..." : "Upload"}
					</button>
				</div>
			)}

			{error && (
				<div className="mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
					<p className="text-red-600 dark:text-red-400">{error}</p>
				</div>
			)}
		</div>
	);
}
