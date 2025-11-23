# Feature Specification: Profile Generation from References

**Feature Branch**: `001-profile-from-references`
**Created**: 2025-11-23
**Status**: Draft
**Input**: User description: "LinkedIn, Indeed are an estabilished standard to show a profile of a person. You can usually see the past work experience and you can also add extra information (f.E. tags about skills or technologies) that highlights key knowledge areas and associates them to the profile and/or to specific work experiences.

Two potential negative aspects of the current standard:

- it's a lot of manual work (so usually many good people have not a lot of information in their profiles and only the people direly searching for a new job put the work in)

- a profile author can claim knowledge, but it doesn't say anything whether they're really profficient in the respective skill

- it doesn't really highlight what past employers actually thought about the employee.

Here's where Credfolio comes in. Instead of manually curating your work experience, you give us your past reference letters and we generate a linkedin like profile for you, with all the tags, achievements, etc. But not only that, we also add what your employer thought about your to add credibility to your profile. You can view this profile in the application. It would look like a LinkedIn profile on steriods, but you could also download a CV directly from the application that's based on the extracted information.



When you've got a job description, we're able to create a custom version of your profile and a custom CV that tweaks how your work experience is described to emphacise a better match towards the job described in the job description.



Credfolio = Portfolio + Credibility"

## Clarifications

### Session 2025-11-23
- Q: How should manual vs. verified data be handled? → A: Hybrid approach; users can edit extracted data and create new unverified entries manually. Visual distinction (e.g., "Verified" vs. "Self-reported") MUST exist for entries backed by reference letters.
- Q: How to handle name mismatches between reference and user profile? → A: Strict Check & Flag; System alerts the user to the mismatch but allows them to proceed by confirming "This is me" (handling marriage, gender changes). The event is logged/flagged for potential future moderation.
- Q: How to handle overlapping data from multiple letters for the same company? → A: Merge by Company; System groups multiple references for the same company (fuzzy match) into a single "Company" entry with aggregated roles/dates to avoid duplication.
- Q: How to handle multiple roles mentioned in a single reference letter? → A: Multi-Role Support; System extracts and stores multiple distinct roles (with separate dates/titles) linked to a single "Company" parent entity to reflect career progression.

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Create Profile from References (Priority: P1)

As a job seeker, I want to upload my past reference letters so that my professional profile is automatically created without manual data entry.

**Why this priority**: This is the core value proposition and entry point. Without this, there is no "Credfolio".

**Independent Test**: Can be tested by uploading a sample reference letter (PDF/DOCX,TXT) and verifying that a profile is created with extracted data.

**Acceptance Scenarios**:

1. **Given** a new user with no profile, **When** they upload a valid reference letter (PDF/DOCX,TXT), **Then** the system parses the file and creates a profile with extracted Name, Role, Company, and Skills.
2. **Given** a user uploading a file, **When** the file format is unsupported (e.g., PNG), **Then** the system shows an error message.
3. **Given** multiple reference letters, **When** uploaded together or sequentially, **Then** the system aggregates the data into a single chronological profile.
4. **Given** a reference letter with a name different from the user's profile name, **When** uploaded, **Then** the system prompts the user to confirm identity before processing, and logs the mismatch flag.
5. **Given** a single reference letter detailing multiple roles (e.g., "started as Junior Dev, promoted to Senior"), **When** processed, **Then** the system creates one Company entry with multiple distinct Role sub-entries.

---

### User Story 2 - View Enhanced Profile (Priority: P1)

As a job seeker, I want to view my generated profile including employer feedback so that I can see how my credibility is presented.

**Why this priority**: Users need to verify and view the value generated from their documents.

**Independent Test**: Can be tested by navigating to the profile view after generation and checking for the presence of "Credibility" or feedback sections.

**Acceptance Scenarios**:

1. **Given** a generated profile, **When** the user views the dashboard, **Then** they see a LinkedIn-style layout with Experience, Skills, and a distinct "Employer Feedback/Credibility" section derived from the letters.
2. **Given** a profile with skills, **When** viewing the Skills section, **Then** skills are linked to the specific work experiences/references where they were mentioned.
3. **Given** a profile with both verified (from letters) and manual entries, **When** viewed, **Then** verified entries are visually distinct (e.g., with a "Verified" badge).
4. **Given** multiple letters for the same company, **When** viewed, **Then** they are displayed as a single Company entry with merged timelines and roles.

---

### User Story 3 - Download Standard CV (Priority: P2)

As a job seeker, I want to download a CV based on my profile so that I can apply to jobs that require a traditional document.

**Why this priority**: Bridges the gap between the platform and traditional hiring processes.

**Independent Test**: Can be tested by clicking "Download CV" and verifying the file content.

**Acceptance Scenarios**:

1. **Given** a complete profile, **When** the user clicks "Download CV", **Then** a PDF file is generated and downloaded containing the profile information formatted as a professional resume.
2. **Given** a profile with "Credibility" data, **When** the CV is generated, **Then** it includes selected quotes or summaries from the reference letters to highlight verified strengths.

---

### User Story 4 - Tailor Profile to Job Description (Priority: P2)

As a job seeker, I want to customize my profile and CV for a specific job description so that I can increase my chances of being a good match.

**Why this priority**: High-value feature that differentiates the product ("LinkedIn on steroids").

**Independent Test**: Can be tested by inputting a sample JD text and verifying that the output changes.

**Acceptance Scenarios**:

1. **Given** a user with a profile and a target Job Description text, **When** they select "Tailor to Job", **Then** the system generates a "Custom Version" of the profile.
2. **Given** a tailored profile version, **When** the user views it, **Then** skills and experiences relevant to the JD are emphasized (e.g., reordered, highlighted).
3. **Given** a tailored profile, **When** the user downloads the CV, **Then** the downloaded PDF reflects the emphasized skills and language of the tailored version.

---

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- What happens when a reference letter is handwritten or scanned (image-based PDF)? (Assume OCR is out of scope for MVP unless specified, or handled by library). -> *System should handle text-based PDFs; scanned images might error or require OCR service (Assume text-based for MVP).*
- What happens when conflicting dates or roles are found in different letters? -> *System should present both or use the most recent/detailed one (Needs strategy, assume mostly additive).*
- How does the system handle languages other than English in reference letters? -> *Assume MVP is English-first.*
- What happens if no skills can be extracted? -> *Profile is created with empty skills section.*
- What happens when a user manually adds an entry? -> *It is marked as "Self-reported" or similar, distinct from "Verified" entries.*
- What happens when the name in the letter doesn't match the user? -> *User is prompted to confirm; system logs a "Name Mismatch" flag for potential review.*
- What happens when multiple letters exist for the same company? -> *They are merged into a single company block with aggregated roles/dates.*
- What happens when one letter contains multiple roles? -> *System parses them as distinct roles under one company.*

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST allow users to upload reference letters in PDF and DOCX formats.
- **FR-002**: System MUST extract structured data from reference letters: Company Name, Job Title, Employment Dates, Candidate Name.
- **FR-003**: System MUST extract "Skills" and "Technologies" keywords from the reference text.
- **FR-004**: System MUST extract subjective "Employer Feedback" (e.g., positive quotes, assessment of strengths) from the reference text.
- **FR-005**: System MUST generate a web-based profile view displaying Experience, Skills, and Credibility/Feedback.
- **FR-006**: System MUST allow users to edit the extracted information (correcting errors or adding missing details).
- **FR-007**: System MUST allow users to manually create new work experience entries that are not backed by a reference letter.
- **FR-008**: System MUST visually distinguish between "Verified" (from reference) and "Self-reported" (manual) entries in the profile view.
- **FR-009**: System MUST generate a downloadable PDF CV based on the profile data.
- **FR-010**: System MUST accept a text input for a "Job Description".
- **FR-011**: System MUST generate a temporary "Tailored" profile/CV version where profile items matching the Job Description keywords are prioritized or highlighted.
- **FR-012**: System MUST detect name mismatches between reference letter and user profile and require user confirmation before processing.
- **FR-013**: System MUST log/flag name mismatch events for future moderation capabilities.
- **FR-014**: System MUST merge extracted data from multiple letters into a single "Company" entry if the company name matches (fuzzy match).
- **FR-015**: System MUST support extracting and storing multiple distinct roles (with different dates/titles) from a single reference letter.

### Key Entities *(include if feature involves data)*

- **UserProfile**: The core entity containing personal info and aggregated professional data.
- **ReferenceLetter**: Represents the source document and its raw extracted metadata.
- **CompanyEntry**: Aggregates work experience at a specific organization. Contains:
  - Company Name
  - **Roles**: List of **WorkExperience** items.
- **WorkExperience**: A specific role at a company. Contains:
  - Role Title
  - Dates
  - Description
  - **Source**: Enum (Verified, Self-reported)
  - **EmployerFeedback** (extracted quotes/sentiment, null for Self-reported)
  - **VerifiedSkills** (skills mentioned in the letter)
  - **MismatchFlag**: Boolean (True if name mismatch was confirmed)
  - **ReferenceIDs**: List of IDs (linking to multiple source letters if merged)
- **JobDescription**: Ephemeral or saved input used for tailoring.
- **GeneratedCV**: The output document (PDF).

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: Users can reach a "Profile Ready" state within 3 minutes of uploading 2 standard reference letters.
- **SC-002**: 90% of generated profiles contain correctly identified Company Names and Job Titles without user correction.
- **SC-003**: Users successfully download a CV in 100% of sessions where a profile was generated.
- **SC-004**: The "Tailored" CV includes at least 20% more keyword overlaps with the Job Description than the generic CV (measured by keyword matching count).
