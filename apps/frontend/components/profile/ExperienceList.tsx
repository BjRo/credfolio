import type { CompanyEntry } from '../../types';
import SkillsList from './SkillsList';

export default function ExperienceList({ companies, highlightSkillIDs }: { companies: CompanyEntry[]; highlightSkillIDs?: string[] }) {
	if (!companies || companies.length === 0) return <div className="text-gray-500 italic">No experience recorded yet.</div>;

	return (
		<div className="space-y-8 relative before:absolute before:inset-0 before:ml-5 before:-translate-x-px md:before:mx-auto md:before:translate-x-0 before:h-full before:w-0.5 before:bg-gradient-to-b before:from-transparent before:via-slate-300 before:to-transparent">
			{companies.map((company) => (
				<div key={company.id} className="relative pl-8 md:pl-0">
					<div className="md:flex items-center justify-between mb-4 md:mb-0 md:flex-row-reverse group">
						{/* Dot */}
						<div className="absolute left-5 md:left-1/2 -translate-x-1/2 w-4 h-4 rounded-full bg-blue-500 border-4 border-white dark:border-gray-900 shadow-sm" />

						<div className="md:w-1/2 md:pl-8 mb-4 md:mb-0">
							{/* Company Info */}
							<h3 className="text-xl font-bold text-gray-800 dark:text-gray-100">{company.name}</h3>
							<p className="text-sm text-gray-500">
								{new Date(company.start_date).toLocaleDateString(undefined, { year: 'numeric', month: 'short' })} -
								{company.end_date ? new Date(company.end_date).toLocaleDateString(undefined, { year: 'numeric', month: 'short' }) : 'Present'}
							</p>
						</div>

						<div className="md:w-1/2 md:pr-8 text-right">
							{/* Placeholder for alignment if needed, or empty */}
						</div>
					</div>

					{/* Roles */}
					<div className="mt-4 space-y-4 md:ml-auto md:w-[calc(50%+2rem)] md:pl-8">
						{company.roles?.map((role) => (
							<div key={role.id} className="bg-white dark:bg-gray-800 p-5 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
								<div className="flex items-center justify-between mb-2">
									<h4 className="font-semibold text-lg text-gray-900 dark:text-white">{role.title}</h4>
									{role.is_verified && (
										<span className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded-full flex items-center gap-1">
											<svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20" role="img" aria-label="Verified">
												<title>Verified</title>
												<path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
											</svg>
											Verified
										</span>
									)}
								</div>
								<p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap mb-4 text-sm leading-relaxed">{role.description}</p>

								{role.employer_feedback && (
									<div className="mt-4 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg text-sm italic text-gray-600 dark:text-gray-400 border-l-4 border-blue-400">
										"{role.employer_feedback}"
									</div>
								)}

								<div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
									<SkillsList skills={role.skills} highlightIDs={highlightSkillIDs} />
								</div>
							</div>
						))}
					</div>
				</div>
			))}
		</div>
	);
}
