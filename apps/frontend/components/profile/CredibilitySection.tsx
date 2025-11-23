export default function CredibilitySection({ companies }: { companies: any[] }) {
	const verifiedRoles = companies?.flatMap(c => c.roles || []).filter((r: any) => r.is_verified) || [];
	const verifiedCount = verifiedRoles.length;
    const totalRoles = companies?.flatMap(c => c.roles || []).length || 0;

    if (verifiedCount === 0) return null;

	return (
		<div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm mb-12 border border-gray-100 dark:border-gray-700">
			<h2 className="text-2xl font-bold mb-6 flex items-center text-gray-900 dark:text-white">
				<span className="mr-2 text-2xl">🛡️</span> Credibility Score
			</h2>

			<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div className="text-center p-6 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-100 dark:border-blue-900/50">
					<div className="text-4xl font-bold text-blue-600 dark:text-blue-400 mb-1">{verifiedCount}</div>
					<div className="text-sm text-gray-600 dark:text-gray-300 font-medium">Verified Roles</div>
				</div>

                <div className="text-center p-6 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-100 dark:border-green-900/50">
					<div className="text-4xl font-bold text-green-600 dark:text-green-400 mb-1">{Math.round((verifiedCount / (totalRoles || 1)) * 100)}%</div>
					<div className="text-sm text-gray-600 dark:text-gray-300 font-medium">Verification Rate</div>
				</div>

                <div className="text-center p-6 bg-purple-50 dark:bg-purple-900/20 rounded-xl border border-purple-100 dark:border-purple-900/50">
					<div className="text-4xl font-bold text-purple-600 dark:text-purple-400 mb-1">A+</div>
					<div className="text-sm text-gray-600 dark:text-gray-300 font-medium">Profile Grade</div>
				</div>
			</div>
		</div>
	);
}

