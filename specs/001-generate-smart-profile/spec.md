# Feature Specification: Generate Smart Profile & Credibility

**Feature Branch**: `001-generate-smart-profile`
**Created**: 2025-11-24
**Status**: Draft
**Input**: User description provided.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Profile Generation from References (Priority: P1)

As a Job Seeker, I want to upload my past reference letters so that my professional profile is automatically populated with my experience, skills, and achievements without manual entry.

**Why this priority**: This is the core value proposition—automating profile creation and verifying it with source documents (credibility).

**Independent Test**: Can be fully tested by uploading a sample reference letter and verifying that the "Experience" and "Skills" sections are populated with correct data.

**Acceptance Scenarios**:

1. **Given** a user with no profile data, **When** they upload a valid PDF reference letter, **Then** the system extracts the company name, role, dates, and key skills.
2. **Given** an uploaded reference letter, **When** processing is complete, **Then** the system populates a "Credibility" section with positive quotes or sentiment summaries from the text.
3. **Given** the generated profile, **When** the user reviews it, **Then** they can edit or delete any extracted information.

---

### User Story 2 - View Credibility Profile (Priority: P1)

As a Job Seeker, I want to view my profile in a "LinkedIn-on-steroids" format that highlights what employers thought of me, so I can verify how my credibility is presented.

**Why this priority**: Essential for the user to validate the "Credibility" aspect before sharing or using the profile.

**Independent Test**: Can be tested by navigating to the profile view after data population.

**Acceptance Scenarios**:

1. **Given** a populated profile, **When** the user views their profile page, **Then** they see standard sections (Experience, Skills) augmented with "Employer Feedback" or "Credibility Highlights".
2. **Given** multiple reference letters processed, **When** viewing the profile, **Then** the system aggregates skills and endorsements across different experiences.

---

### User Story 3 - Tailor Profile to Job Description (Priority: P2)

As a Job Seeker, I want to provide a job description for a role I'm interested in, so that the system generates a custom version of my profile and CV that emphasizes the most relevant experience.

**Why this priority**: High value add for applying to specific jobs, but depends on the base profile existing (P1).

**Independent Test**: Upload a JD and compare the standard profile vs. the tailored profile.

**Acceptance Scenarios**:

1. **Given** an existing profile and a Job Description text, **When** the user requests a "Tailored Match", **Then** the system reorders or highlights skills/experiences that match the JD keywords.
2. **Given** a tailored profile, **When** the user views it, **Then** they see a "Match Score" or explanation of why certain elements are highlighted.

---

### User Story 4 - Download CV (Priority: P2)

As a Job Seeker, I want to download my standard or tailored profile as a formatted CV, so I can submit it to application portals.

**Why this priority**: Crucial for the "end of loop" utility (actually applying to jobs).

**Independent Test**: detailed generation of PDF file.

**Acceptance Scenarios**:

1. **Given** a standard or tailored profile, **When** the user clicks "Download CV", **Then** a PDF file is generated and downloaded.
2. **Given** a tailored profile, **When** downloaded, **Then** the PDF reflects the emphasized content specific to that job description.

---

### Edge Cases

- What happens when a reference letter is handwritten or low-quality scan? (Assume OCR failure handling or rejection).
- How does the system handle conflicting dates or roles from different letters? (User should resolve, or show all).
- What happens if the Job Description is too short or vague? (System provides best-effort match or asks for more info).
- What if the reference letter is in a language other than English? (Assume MVP supports English only or standard major languages).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to upload reference letters in PDF, DOCX, and TXT formats.
- **FR-002**: System MUST extract structured data from reference letters: Company Name, Job Title, Employment Dates, Skills, and Achievements.
- **FR-003**: System MUST extract "Credibility Indicators" (positive sentiment, specific praise quotes) from reference letters.
- **FR-004**: System MUST provide an interface for users to review, edit, and approve extracted data.
- **FR-005**: System MUST display a comprehensive profile view combining manual edits and extracted data.
- **FR-006**: System MUST allow users to input a Job Description (paste text or upload file).
- **FR-007**: System MUST generate a "Tailored Profile" version that ranks experience/skills based on semantic similarity to the Job Description.
- **FR-008**: System MUST generate a downloadable PDF CV from the current profile view (Standard or Tailored).
- **FR-009**: System MUST persist user profiles and uploaded documents securely.

### Key Entities

- **UserProfile**: Aggregates all career info (User ID, Summary, Skills List).
- **WorkExperience**: Derived from reference letters (Company, Role, Dates, Description, ReferenceSourceID).
- **ReferenceLetter**: The source document (File Path, Upload Date, Parsed Status).
- **CredibilityHighlight**: Specific quote or sentiment tag linked to a WorkExperience.
- **JobMatch**: A session or record of a specific JD vs. Profile comparison (TargetRole, MatchAnalysis, TailoredCV_Path).

### Assumptions

- Users already have accounts and are authenticated.
- Reference letters are in supported languages (initially English).
- The extraction technology (OCR/AI) is available as a service or module.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can go from uploading a reference letter to viewing a populated draft profile in under 60 seconds.
- **SC-002**: 90% of users utilize the "Download CV" feature after generating a profile (indicating value).
- **SC-003**: The system successfully extracts Company and Role from 95% of standard typed reference letters.
- **SC-004**: Users make edits to less than 20% of the auto-populated fields (indicating high extraction accuracy).
