"use client";

import { useState } from "react";
import { uploadReferenceLetter } from "../../lib/api/referenceLetters";
import { getErrorMessage } from "../../lib/utils/errorMessages";
import LoadingSpinner from "./LoadingSpinner";

export default function ReferenceLetterUpload({
  onUploadComplete,
}: {
  onUploadComplete?: () => void;
}) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    setError(null);
    setSuccess(null);

    try {
      await uploadReferenceLetter(file);
      setSuccess("Reference letter uploaded successfully!");
      if (onUploadComplete) onUploadComplete();
      // Reset input - might be tricky with React state, but we can leave it populated or clear it via ref if needed.
      // For simple UI, just leaving it is fine or user can upload another.
    } catch (err) {
      setError(getErrorMessage(err));
      console.error(err);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="p-4 border rounded-lg shadow-sm bg-white">
      <h3 className="text-lg font-medium mb-2">Upload Reference Letter</h3>
      <p className="text-sm text-gray-500 mb-4">
        Upload a reference letter to generate your profile.
      </p>

      <div className="flex items-center gap-4">
        <label className="block flex-1">
          <span className="sr-only">Choose file</span>
          <input
            type="file"
            accept=".md, .txt"
            onChange={handleFileChange}
            disabled={uploading}
            className="block w-full text-sm text-slate-500
              file:mr-4 file:py-2 file:px-4
              file:rounded-full file:border-0
              file:text-sm file:font-semibold
              file:bg-violet-50 file:text-violet-700
              hover:file:bg-violet-100
              disabled:opacity-50"
          />
        </label>
        {uploading && (
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <LoadingSpinner size="sm" className="text-indigo-600" />
            <span>Uploading and processing...</span>
          </div>
        )}
      </div>

      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
      {success && <p className="mt-2 text-sm text-green-600">{success}</p>}
    </div>
  );
}
