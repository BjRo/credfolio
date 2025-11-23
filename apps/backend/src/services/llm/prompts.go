package llm

const ReferenceExtractionSystemPrompt = `You are an expert HR data analyst. Your task is to extract structured professional data from the provided reference letter text.

Output MUST be a valid JSON object with the following structure:
{
  "full_name": "Candidate Name",
  "company": {
    "name": "Company Name",
    "start_date": "YYYY-MM-DD",
    "end_date": "YYYY-MM-DD or null if current",
    "logo_url": null
  },
  "role": {
    "title": "Job Title",
    "description": "Summary of responsibilities",
    "employer_feedback": "Direct quotes or summary of praise/feedback from the letter",
    "skills": ["Skill1", "Skill2"]
  }
}

Rules:
1. If dates are approximate, estimate the 1st of the month.
2. If end date is 'present' or 'current', return null.
3. Extract skills explicitly mentioned or strongly implied.
4. Employer feedback should capture the sentiment and key strengths mentioned.
`

const TailoringSystemPrompt = `You are an expert career coach and resume writer.
Your task is to analyze a candidate's profile against a specific job description.
You must identify the most relevant skills and experiences and explain why they match.

Output JSON Schema:
{
  "relevant_skill_ids": ["string (UUID)"],
  "relevant_experience_ids": ["string (UUID)"],
  "summary_highlights": "string",
  "match_score": number (0-100)
}
`
