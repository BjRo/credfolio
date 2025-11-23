export default function ProfileHeader({ profile }: { profile: any }) {
	if (!profile) return null;

	const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
	const cvUrl = `${apiUrl}/api/profile/cv?user_id=${profile.id}`;

	return (
		<div className="mb-12 text-center">
			<div className="inline-block p-1 rounded-full bg-gradient-to-r from-blue-500 to-teal-400 mb-4">
				<div className="w-24 h-24 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-3xl font-bold text-gray-400">
					{profile.full_name?.charAt(0) || '?'}
				</div>
			</div>
			<h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">{profile.full_name}</h1>
			<p className="text-gray-500 dark:text-gray-400 text-lg mb-6">{profile.email}</p>

			<a
				href={cvUrl}
				download="credfolio_cv.pdf"
				className="inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 md:py-4 md:text-lg shadow-sm transition-colors"
			>
				<span className="mr-2">📄</span> Download CV
			</a>
		</div>
	);
}
