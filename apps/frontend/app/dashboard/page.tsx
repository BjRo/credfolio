import { getProfile } from "../../services/api";
import type { UserProfile } from "../../types";
import DashboardView from "../../components/profile/DashboardView";

export default async function DashboardPage({
	searchParams,
}: {
	searchParams: { [key: string]: string | string[] | undefined };
}) {
	const userID = searchParams.user_id as string;

	if (!userID) {
		return (
			<div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 p-4">
				<div className="text-center p-8 bg-white dark:bg-gray-800 rounded-xl shadow-lg">
					<h2 className="text-xl font-bold text-gray-900 dark:text-white mb-2">
						User ID Required
					</h2>
					<p className="text-gray-500 mb-4">
						Please provide a user_id in the URL query parameters.
					</p>
					<code className="bg-gray-100 dark:bg-gray-900 px-2 py-1 rounded text-sm">
						?user_id=...
					</code>
				</div>
			</div>
		);
	}

	let profile: UserProfile;
	try {
		profile = await getProfile(userID);
	} catch (e) {
		return (
			<div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 p-4">
				<div className="text-center p-8 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-red-100">
					<h2 className="text-xl font-bold text-red-600 mb-2">
						Error Loading Profile
					</h2>
					<p className="text-gray-500">{String(e)}</p>
				</div>
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
			<DashboardView profile={profile} userID={userID} />
		</div>
	);
}
