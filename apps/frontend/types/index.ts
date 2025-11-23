export interface Skill {
	id: string;
	name: string;
}

export interface WorkExperience {
	id: string;
	company_id: string;
	title: string;
	start_date: string;
	end_date?: string | null;
	description: string;
	source: string;
	employer_feedback: string;
	reference_letter_id?: string | null;
	is_verified: boolean;
	skills?: Skill[];
}

export interface CompanyEntry {
	id: string;
	user_id: string;
	name: string;
	logo_url?: string | null;
	start_date: string;
	end_date?: string | null;
	roles?: WorkExperience[];
}

export interface UserProfile {
	id: string;
	email: string;
	full_name: string;
	created_at: string;
	companies?: CompanyEntry[];
}

export interface TailoringResult {
	relevant_skill_ids: string[];
	relevant_experience_ids: string[];
	summary_highlights: string;
	match_score: number;
}

