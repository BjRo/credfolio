import FileUploader from '../../components/upload/FileUploader';

export default function UploadPage() {
	return (
		<div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 p-4">
			<h1 className="text-3xl font-bold mb-8 text-gray-900 dark:text-white">
				Upload Reference Letter
			</h1>
			<p className="text-gray-600 dark:text-gray-300 mb-8 text-center max-w-lg">
				Upload your reference letter (PDF) to automatically generate your Credfolio profile.
			</p>
			<FileUploader />
		</div>
	);
}

