"use client";

import { useState } from 'react';
import { uploadReferenceLetter } from '../../services/api';

export default function FileUploader() {
	const [status, setStatus] = useState<'idle' | 'uploading' | 'success' | 'error'>('idle');
	const [message, setMessage] = useState('');

	const handleFile = async (file: File) => {
		if (file.type !== 'application/pdf') {
			setStatus('error');
			setMessage('Only PDF files are supported');
			return;
		}

		setStatus('uploading');
		try {
			await uploadReferenceLetter(file);
			setStatus('success');
			setMessage('File uploaded successfully! Processing...');
		} catch (err) {
			setStatus('error');
			setMessage('Upload failed. Please try again.');
		}
	};

	const onDrop = (e: React.DragEvent) => {
		e.preventDefault();
		if (e.dataTransfer.files && e.dataTransfer.files[0]) {
			handleFile(e.dataTransfer.files[0]);
		}
	};

	const onDragOver = (e: React.DragEvent) => {
		e.preventDefault();
	};

	const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		if (e.target.files && e.target.files[0]) {
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
					accept=".pdf"
					className="hidden"
					id="file-upload"
					onChange={onChange}
				/>
				<label htmlFor="file-upload" className="cursor-pointer block w-full h-full">
					<div className="text-lg mb-2 font-medium">Drag & Drop or Click to Upload</div>
					<div className="text-sm text-gray-500">PDF files only</div>
				</label>
			</div>

			{status === 'uploading' && (
				<div className="mt-4 text-blue-600 text-center animate-pulse">Uploading...</div>
			)}
			{status === 'success' && (
				<div className="mt-4 text-green-600 text-center">{message}</div>
			)}
			{status === 'error' && (
				<div className="mt-4 text-red-600 text-center">{message}</div>
			)}
		</div>
	);
}

