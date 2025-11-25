---
description: "Task list for Generate Smart Profile & Credibility feature implementation"
---

# Tasks: Generate Smart Profile & Credibility

**Input**: Design documents from `/specs/001-generate-smart-profile/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Unit tests are MANDATORY for all developments per the Unit-Testing-First principle. All code changes MUST include corresponding unit tests. Integration tests are optional and may be included when explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Monorepo**: `apps/backend/`, `apps/frontend/`
- **Backend**: `apps/backend/internal/{domain,service,repository,handler}`, `apps/backend/pkg/{ai,pdf,extraction}`
- **Frontend**: `apps/frontend/app/profile`, `apps/frontend/components/profile`, `apps/frontend/lib/api`
- **Backend Tests**: `apps/backend/internal/{domain,service,repository,handler}/..._test.go`
- **Frontend Tests**: `apps/frontend/components/profile/*.test.tsx`, `apps/frontend/lib/api/*.test.ts`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create backend project structure per implementation plan in apps/backend/internal/domain, apps/backend/internal/service, apps/backend/internal/repository, apps/backend/internal/handler
- [ ] T002 Create frontend project structure per implementation plan in apps/frontend/app/profile, apps/frontend/components/profile, apps/frontend/lib/api
- [ ] T003 [P] Initialize Go dependencies in apps/backend/go.mod: gorm.io/gorm, gorm.io/driver/postgres, github.com/go-chi/chi, github.com/openai/openai-go, github.com/johnfercher/maroto/v2, github.com/ledongthuc/pdf, github.com/oapi-codegen/oapi-codegen/v2
- [ ] T004 [P] Configure backend linting and formatting tools (golangci-lint, gofmt) in apps/backend/
- [ ] T005 [P] Configure frontend linting and formatting tools (biome, prettier) in apps/frontend/
- [ ] T006 Setup environment configuration management in apps/backend/pkg/config/config.go
- [ ] T007 [P] Create OpenAPI contract file in apps/backend/api/openapi.yaml from specs/001-generate-smart-profile/contracts/openapi.yaml

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Unit Tests for Foundational Components (MANDATORY) ⚠️

> **NOTE: Write unit tests FIRST, ensure they FAIL before implementation. Follow AAA style, hide implementation details, use descriptive names (context_trigger_expectation), and stub external calls.**

- [ ] T008 [P] Unit test for User model validation in apps/backend/internal/domain/user_test.go
- [ ] T009 [P] Unit test for Profile model validation in apps/backend/internal/domain/profile_test.go
- [ ] T010 [P] Unit test for WorkExperience model validation in apps/backend/internal/domain/work_experience_test.go
- [ ] T011 [P] Unit test for Skill model validation in apps/backend/internal/domain/skill_test.go
- [ ] T012 [P] Unit test for ReferenceLetter model validation in apps/backend/internal/domain/reference_letter_test.go
- [ ] T013 [P] Unit test for CredibilityHighlight model validation in apps/backend/internal/domain/credibility_highlight_test.go
- [ ] T014 [P] Unit test for JobMatch model validation in apps/backend/internal/domain/job_match_test.go
- [ ] T015 [P] Unit test for LLMProvider interface mock in apps/backend/internal/service/llm_provider_test.go
- [ ] T016 [P] Unit test for OpenAIProvider when given valid prompt returns structured data in apps/backend/pkg/ai/openai_provider_test.go
- [ ] T017 [P] Unit test for OpenAIProvider when API call fails returns error in apps/backend/pkg/ai/openai_provider_test.go
- [ ] T018 [P] Unit test for text extractor when given txt file returns text content in apps/backend/pkg/extraction/extractor_test.go
- [ ] T019 [P] Unit test for text extractor when given markdown file returns text content in apps/backend/pkg/extraction/extractor_test.go
- [ ] T020 [P] Unit test for PDF generator when given profile data generates PDF bytes in apps/backend/pkg/pdf/generator_test.go
- [ ] T021 [P] Unit test for config loader when given valid env vars loads configuration in apps/backend/pkg/config/config_test.go
- [ ] T022 [P] Unit test for logger when logging info message writes to output in apps/backend/pkg/logger/logger_test.go
- [ ] T023 [P] Unit test for auth middleware when request has no user injects mock user in apps/backend/internal/handler/middleware/auth_test.go

### Implementation for Foundational Components

- [ ] T024 Setup PostgreSQL database connection using GORM in apps/backend/internal/repository/db.go
- [ ] T025 Create base User model (mock/existing) in apps/backend/internal/domain/user.go
- [ ] T026 [P] Create base Profile model in apps/backend/internal/domain/profile.go
- [ ] T027 [P] Create base WorkExperience model in apps/backend/internal/domain/work_experience.go
- [ ] T028 [P] Create base Skill model in apps/backend/internal/domain/skill.go
- [ ] T029 [P] Create base ReferenceLetter model in apps/backend/internal/domain/reference_letter.go
- [ ] T030 [P] Create base CredibilityHighlight model in apps/backend/internal/domain/credibility_highlight.go
- [ ] T031 [P] Create base JobMatch model in apps/backend/internal/domain/job_match.go
- [ ] T032 Setup GORM AutoMigrate for all domain models in apps/backend/internal/repository/migrations.go
- [ ] T033 [P] Create LLMProvider interface abstraction in apps/backend/internal/service/llm_provider.go
- [ ] T034 [P] Implement OpenAIProvider struct satisfying LLMProvider interface in apps/backend/pkg/ai/openai_provider.go
- [ ] T035 [P] Create text extraction service in apps/backend/pkg/extraction/extractor.go
- [ ] T036 [P] Create PDF generation service using maroto in apps/backend/pkg/pdf/generator.go
- [ ] T037 Setup Chi router and middleware structure in apps/backend/cmd/server/main.go
- [ ] T038 Create mock authentication middleware (inject current user ID) in apps/backend/internal/handler/middleware/auth.go
- [ ] T039 Configure error handling and logging infrastructure in apps/backend/pkg/logger/logger.go
- [ ] T040 Setup environment configuration management in apps/backend/pkg/config/config.go
- [ ] T041 Generate OpenAPI server stubs using oapi-codegen in apps/backend/api/generated/
- [ ] T042 Generate TypeScript API client from OpenAPI spec in apps/frontend/lib/api/generated/

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Profile Generation from References (Priority: P1) 🎯 MVP

**Goal**: Automatically populate professional profile by extracting data from uploaded reference letters using AI, including credibility highlights from employer sentiment.

**Independent Test**: Upload a sample reference letter and verify that the "Experience" and "Skills" sections are populated with correct data, and that "Credibility" section contains positive quotes or sentiment summaries.

### Unit Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write unit tests FIRST, ensure they FAIL before implementation. Follow AAA style, hide implementation details, use descriptive names (context_trigger_expectation), and stub external calls.**

- [ ] T043 [P] [US1] Unit test for ProfileRepository when saving profile persists to database in apps/backend/internal/repository/profile_repository_test.go
- [ ] T044 [P] [US1] Unit test for ProfileRepository when finding by user ID returns profile in apps/backend/internal/repository/profile_repository_test.go
- [ ] T045 [P] [US1] Unit test for ReferenceLetterRepository when saving reference letter persists to database in apps/backend/internal/repository/reference_letter_repository_test.go
- [ ] T046 [P] [US1] Unit test for ReferenceLetterRepository when finding by user ID returns letters in apps/backend/internal/repository/reference_letter_repository_test.go
- [ ] T047 [P] [US1] Unit test for WorkExperienceRepository when saving work experience persists to database in apps/backend/internal/repository/work_experience_repository_test.go
- [ ] T048 [P] [US1] Unit test for CredibilityHighlightRepository when saving highlight persists to database in apps/backend/internal/repository/credibility_highlight_repository_test.go
- [ ] T049 [P] [US1] Unit test for ProfileService when generating profile from reference letter extracts structured data in apps/backend/internal/service/profile_service_test.go
- [ ] T050 [P] [US1] Unit test for ProfileService when LLMProvider returns error propagates error in apps/backend/internal/service/profile_service_test.go
- [ ] T051 [P] [US1] Unit test for ProfileService when extracting credibility finds positive sentiment quotes in apps/backend/internal/service/profile_service_test.go
- [ ] T052 [P] [US1] Unit test for ProfileService when reference letter is invalid returns validation error in apps/backend/internal/service/profile_service_test.go
- [ ] T053 [P] [US1] Unit test for ReferenceLetterHandler when uploading valid file returns reference letter ID in apps/backend/internal/handler/reference_letter_handler_test.go
- [ ] T054 [P] [US1] Unit test for ReferenceLetterHandler when uploading invalid file type returns error in apps/backend/internal/handler/reference_letter_handler_test.go
- [ ] T055 [P] [US1] Unit test for ProfileHandler when generating profile returns profile data in apps/backend/internal/handler/profile_handler_test.go
- [ ] T056 [P] [US1] Unit test for ReferenceLetterUpload component when file selected shows file name in apps/frontend/components/profile/ReferenceLetterUpload.test.tsx
- [ ] T057 [P] [US1] Unit test for ReferenceLetterUpload component when upload fails displays error message in apps/frontend/components/profile/ReferenceLetterUpload.test.tsx
- [ ] T058 [P] [US1] Unit test for GenerateProfileButton component when clicked triggers generation in apps/frontend/components/profile/GenerateProfileButton.test.tsx
- [ ] T059 [P] [US1] Unit test for ProfileEditor component when editing field updates value in apps/frontend/components/profile/ProfileEditor.test.tsx
- [ ] T060 [P] [US1] Unit test for API client when posting reference letter sends multipart form data in apps/frontend/lib/api/referenceLetters.test.ts
- [ ] T061 [P] [US1] Unit test for API client when generating profile calls correct endpoint in apps/frontend/lib/api/profile.test.ts

### Implementation for User Story 1

- [ ] T062 [P] [US1] Create Profile repository interface in apps/backend/internal/repository/profile_repository.go
- [ ] T063 [P] [US1] Create ReferenceLetter repository interface in apps/backend/internal/repository/reference_letter_repository.go
- [ ] T064 [P] [US1] Create WorkExperience repository interface in apps/backend/internal/repository/work_experience_repository.go
- [ ] T065 [P] [US1] Create CredibilityHighlight repository interface in apps/backend/internal/repository/credibility_highlight_repository.go
- [ ] T066 [P] [US1] Implement GormProfileRepository in apps/backend/internal/repository/gorm_profile_repository.go
- [ ] T067 [P] [US1] Implement GormReferenceLetterRepository in apps/backend/internal/repository/gorm_reference_letter_repository.go
- [ ] T068 [P] [US1] Implement GormWorkExperienceRepository in apps/backend/internal/repository/gorm_work_experience_repository.go
- [ ] T069 [P] [US1] Implement GormCredibilityHighlightRepository in apps/backend/internal/repository/gorm_credibility_highlight_repository.go
- [ ] T070 [US1] Create ProfileService with GenerateProfileFromReferences method in apps/backend/internal/service/profile_service.go
- [ ] T071 [US1] Implement AI extraction logic in ProfileService that calls LLMProvider to extract structured data (Company, Role, Dates, Skills, Achievements) in apps/backend/internal/service/profile_service.go
- [ ] T072 [US1] Implement credibility extraction logic that extracts positive sentiment quotes from reference letters in apps/backend/internal/service/profile_service.go
- [ ] T073 [US1] Implement file upload handler for reference letters (multipart/form-data) in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T074 [US1] Implement POST /reference-letters endpoint handler in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T075 [US1] Implement POST /profile/generate endpoint handler that processes uploaded reference letters in apps/backend/internal/handler/profile_handler.go
- [ ] T076 [US1] Add validation and error handling for file uploads and AI extraction in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T077 [US1] Add logging for profile generation operations in apps/backend/internal/service/profile_service.go
- [ ] T078 [US1] Create reference letter upload UI component in apps/frontend/components/profile/ReferenceLetterUpload.tsx
- [ ] T079 [US1] Create profile generation trigger UI in apps/frontend/components/profile/GenerateProfileButton.tsx
- [ ] T080 [US1] Create API client method for POST /reference-letters in apps/frontend/lib/api/referenceLetters.ts
- [ ] T081 [US1] Create API client method for POST /profile/generate in apps/frontend/lib/api/profile.ts
- [ ] T082 [US1] Create profile generation page in apps/frontend/app/profile/generate/page.tsx
- [ ] T083 [US1] Implement profile data editing interface (edit/delete extracted information) in apps/frontend/components/profile/ProfileEditor.tsx

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - users can upload reference letters, generate profiles, and edit extracted data.

---

## Phase 4: User Story 2 - View Credibility Profile (Priority: P1)

**Goal**: Display profile in a comprehensive "LinkedIn-on-steroids" format that highlights employer feedback and credibility highlights, with aggregated skills and endorsements across experiences.

**Independent Test**: Navigate to the profile view after data population and verify standard sections (Experience, Skills) are augmented with "Employer Feedback" or "Credibility Highlights", and that skills/endorsements are aggregated across multiple reference letters.

### Unit Tests for User Story 2 (MANDATORY) ⚠️

> **NOTE: Write unit tests FIRST, ensure they FAIL before implementation. Follow AAA style, hide implementation details, use descriptive names (context_trigger_expectation), and stub external calls.**

- [ ] T084 [P] [US2] Unit test for ProfileService when getting profile aggregates skills across experiences in apps/backend/internal/service/profile_service_test.go
- [ ] T085 [P] [US2] Unit test for ProfileService when getting profile includes credibility highlights in apps/backend/internal/service/profile_service_test.go
- [ ] T086 [P] [US2] Unit test for ProfileHandler when getting profile returns complete profile data in apps/backend/internal/handler/profile_handler_test.go
- [ ] T087 [P] [US2] Unit test for ProfileHandler when profile not found returns 404 error in apps/backend/internal/handler/profile_handler_test.go
- [ ] T088 [P] [US2] Unit test for ProfileView component when profile loaded displays all sections in apps/frontend/components/profile/ProfileView.test.tsx
- [ ] T089 [P] [US2] Unit test for ProfileView component when loading shows loading state in apps/frontend/components/profile/ProfileView.test.tsx
- [ ] T090 [P] [US2] Unit test for CredibilityHighlights component when given highlights displays quotes in apps/frontend/components/profile/CredibilityHighlights.test.tsx
- [ ] T091 [P] [US2] Unit test for WorkExperienceCard component when given experience displays credibility highlights in apps/frontend/components/profile/WorkExperienceCard.test.tsx
- [ ] T092 [P] [US2] Unit test for SkillsSection component when given skills displays aggregated list in apps/frontend/components/profile/SkillsSection.test.tsx
- [ ] T093 [P] [US2] Unit test for API client when getting profile returns profile data in apps/frontend/lib/api/profile.test.ts

### Implementation for User Story 2

- [ ] T094 [US2] Implement GET /profile endpoint handler in apps/backend/internal/handler/profile_handler.go
- [ ] T095 [US2] Add logic to aggregate skills and endorsements across multiple work experiences in ProfileService.GetProfile method in apps/backend/internal/service/profile_service.go
- [ ] T096 [US2] Create API client method for GET /profile in apps/frontend/lib/api/profile.ts
- [ ] T097 [US2] Create ProfileView component displaying Experience, Skills, and Credibility Highlights in apps/frontend/components/profile/ProfileView.tsx
- [ ] T098 [US2] Create CredibilityHighlights section component in apps/frontend/components/profile/CredibilityHighlights.tsx
- [ ] T099 [US2] Create WorkExperience display component with credibility highlights in apps/frontend/components/profile/WorkExperienceCard.tsx
- [ ] T100 [US2] Create Skills aggregation display component in apps/frontend/components/profile/SkillsSection.tsx
- [ ] T101 [US2] Create profile view page in apps/frontend/app/profile/page.tsx
- [ ] T102 [US2] Add styling for "LinkedIn-on-steroids" profile layout using Tailwind CSS in apps/frontend/components/profile/ProfileView.tsx

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - users can generate profiles and view them with credibility highlights.

---

## Phase 5: User Story 3 - Tailor Profile to Job Description (Priority: P2)

**Goal**: Generate a custom version of profile and CV that emphasizes the most relevant experience and skills based on semantic similarity to a provided job description, with match score and explanation.

**Independent Test**: Upload a Job Description and compare the standard profile vs. the tailored profile, verifying that skills/experiences matching JD keywords are reordered or highlighted, and that a Match Score or explanation is displayed.

### Unit Tests for User Story 3 (MANDATORY) ⚠️

> **NOTE: Write unit tests FIRST, ensure they FAIL before implementation. Follow AAA style, hide implementation details, use descriptive names (context_trigger_expectation), and stub external calls.**

- [ ] T103 [P] [US3] Unit test for JobMatchRepository when saving job match persists to database in apps/backend/internal/repository/job_match_repository_test.go
- [ ] T104 [P] [US3] Unit test for JobMatchRepository when finding by profile ID returns matches in apps/backend/internal/repository/job_match_repository_test.go
- [ ] T105 [P] [US3] Unit test for TailoringService when tailoring profile ranks experiences by relevance in apps/backend/internal/service/tailoring_service_test.go
- [ ] T106 [P] [US3] Unit test for TailoringService when job description is empty returns error in apps/backend/internal/service/tailoring_service_test.go
- [ ] T107 [P] [US3] Unit test for TailoringService when calculating match score returns score between 0 and 1 in apps/backend/internal/service/tailoring_service_test.go
- [ ] T108 [P] [US3] Unit test for TailoringService when LLMProvider fails propagates error in apps/backend/internal/service/tailoring_service_test.go
- [ ] T109 [P] [US3] Unit test for ProfileHandler when tailoring profile returns tailored profile data in apps/backend/internal/handler/profile_handler_test.go
- [ ] T110 [P] [US3] Unit test for ProfileHandler when job description invalid returns validation error in apps/backend/internal/handler/profile_handler_test.go
- [ ] T111 [P] [US3] Unit test for JobDescriptionInput component when text entered updates value in apps/frontend/components/profile/JobDescriptionInput.test.tsx
- [ ] T112 [P] [US3] Unit test for TailoredProfileView component when given tailored profile highlights matched content in apps/frontend/components/profile/TailoredProfileView.test.tsx
- [ ] T113 [P] [US3] Unit test for MatchScore component when given score displays percentage in apps/frontend/components/profile/MatchScore.test.tsx
- [ ] T114 [P] [US3] Unit test for API client when tailoring profile sends job description in apps/frontend/lib/api/profile.test.ts

### Implementation for User Story 3

- [ ] T115 [US3] Create JobMatch repository interface in apps/backend/internal/repository/job_match_repository.go
- [ ] T116 [US3] Implement GormJobMatchRepository in apps/backend/internal/repository/gorm_job_match_repository.go
- [ ] T117 [US3] Create TailoringService with TailorProfileToJobDescription method in apps/backend/internal/service/tailoring_service.go
- [ ] T118 [US3] Implement semantic matching logic using LLMProvider to rank experience/skills based on job description in apps/backend/internal/service/tailoring_service.go
- [ ] T119 [US3] Implement match score calculation in TailoringService in apps/backend/internal/service/tailoring_service.go
- [ ] T120 [US3] Implement POST /profile/tailor endpoint handler in apps/backend/internal/handler/profile_handler.go
- [ ] T121 [US3] Add validation for job description input in apps/backend/internal/handler/profile_handler.go
- [ ] T122 [US3] Create API client method for POST /profile/tailor in apps/frontend/lib/api/profile.ts
- [ ] T123 [US3] Create JobDescriptionInput component in apps/frontend/components/profile/JobDescriptionInput.tsx
- [ ] T124 [US3] Create TailoredProfileView component showing highlighted/reordered content in apps/frontend/components/profile/TailoredProfileView.tsx
- [ ] T125 [US3] Create MatchScore display component in apps/frontend/components/profile/MatchScore.tsx
- [ ] T126 [US3] Create profile tailoring page in apps/frontend/app/profile/tailor/page.tsx
- [ ] T127 [US3] Add explanation UI for why certain elements are highlighted in apps/frontend/components/profile/TailoredProfileView.tsx

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently - users can generate profiles, view them, and tailor them to job descriptions.

---

## Phase 6: User Story 4 - Download CV (Priority: P2)

**Goal**: Generate and download a formatted PDF CV from the current profile view (Standard or Tailored), reflecting emphasized content when tailored.

**Independent Test**: Generate a PDF file from a standard or tailored profile and verify the PDF contains all profile data in a formatted CV layout, with tailored content emphasized when applicable.

### Unit Tests for User Story 4 (MANDATORY) ⚠️

> **NOTE: Write unit tests FIRST, ensure they FAIL before implementation. Follow AAA style, hide implementation details, use descriptive names (context_trigger_expectation), and stub external calls.**

- [ ] T128 [P] [US4] Unit test for CV generator when given profile data generates PDF with all sections in apps/backend/pkg/pdf/cv_generator_test.go
- [ ] T129 [P] [US4] Unit test for CV generator when given tailored profile emphasizes matched content in apps/backend/pkg/pdf/cv_generator_test.go
- [ ] T130 [P] [US4] Unit test for CV generator when profile is empty returns error in apps/backend/pkg/pdf/cv_generator_test.go
- [ ] T131 [P] [US4] Unit test for ProfileHandler when downloading CV returns PDF bytes in apps/backend/internal/handler/profile_handler_test.go
- [ ] T132 [P] [US4] Unit test for ProfileHandler when profile not found returns 404 error in apps/backend/internal/handler/profile_handler_test.go
- [ ] T133 [P] [US4] Unit test for ProfileHandler when PDF generation fails returns error in apps/backend/internal/handler/profile_handler_test.go
- [ ] T134 [P] [US4] Unit test for DownloadCVButton component when clicked triggers download in apps/frontend/components/profile/DownloadCVButton.test.tsx
- [ ] T135 [P] [US4] Unit test for DownloadCVButton component when download fails displays error in apps/frontend/components/profile/DownloadCVButton.test.tsx
- [ ] T136 [P] [US4] Unit test for API client when downloading CV returns blob data in apps/frontend/lib/api/profile.test.ts

### Implementation for User Story 4

- [ ] T137 [US4] Enhance PDF generation service to create CV layout using maroto in apps/backend/pkg/pdf/cv_generator.go
- [ ] T138 [US4] Implement CV template with sections: Summary, Work Experience, Skills, Credibility Highlights in apps/backend/pkg/pdf/cv_generator.go
- [ ] T139 [US4] Add logic to emphasize tailored content in PDF when JobMatch is provided in apps/backend/pkg/pdf/cv_generator.go
- [ ] T140 [US4] Implement GET /profile/{profileId}/cv endpoint handler in apps/backend/internal/handler/profile_handler.go
- [ ] T141 [US4] Add query parameter support for tailored CV (jobMatchId) in GET /profile/{profileId}/cv handler in apps/backend/internal/handler/profile_handler.go
- [ ] T142 [US4] Create API client method for GET /profile/{profileId}/cv in apps/frontend/lib/api/profile.ts
- [ ] T143 [US4] Create DownloadCVButton component in apps/frontend/components/profile/DownloadCVButton.tsx
- [ ] T144 [US4] Add download CV functionality to profile view page in apps/frontend/app/profile/page.tsx
- [ ] T145 [US4] Add download CV functionality to tailored profile view in apps/frontend/app/profile/tailor/page.tsx
- [ ] T146 [US4] Add error handling for PDF generation failures in apps/backend/internal/handler/profile_handler.go

**Checkpoint**: All user stories should now be independently functional - users can generate profiles, view them with credibility, tailor them to jobs, and download CVs.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T147 [P] Documentation updates in apps/backend/api/openapi.yaml with complete request/response examples
- [ ] T148 [P] Add comprehensive error messages and user feedback throughout frontend components
- [ ] T149 Code cleanup and refactoring: Extract common patterns in service layer
- [ ] T150 Performance optimization: Cache AI responses where appropriate in apps/backend/internal/service/profile_service.go
- [ ] T151 [P] Add loading states and progress indicators for AI operations in apps/frontend/components/profile/
- [ ] T152 [P] Add input validation and sanitization for all user inputs in apps/backend/internal/handler/
- [ ] T153 Security hardening: Validate file types and sizes for reference letter uploads in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T154 Run quickstart.md validation: Verify all setup steps work end-to-end
- [ ] T155 [P] Add comprehensive logging for all API endpoints in apps/backend/internal/handler/
- [ ] T156 Optimize database queries with proper indexing in apps/backend/internal/repository/
- [ ] T157 Add rate limiting for AI API calls in apps/backend/internal/service/profile_service.go
- [ ] T158 [P] Additional unit tests for edge cases in apps/backend/internal/service/ and apps/frontend/components/profile/

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 for profile data structure, but should be independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Depends on US1/US2 for base profile, but should be independently testable
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Depends on US1/US2 for profile data, and US3 for tailored content, but should be independently testable

### Within Each User Story

- Unit tests MUST be written FIRST and FAIL before implementation
- Models before repositories
- Repositories before services
- Services before handlers/endpoints
- Backend API before frontend integration
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, User Stories 1 and 2 (both P1) can start in parallel (if team capacity allows)
- All unit tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Repository implementations marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members (after foundational phase)

---

## Parallel Example: User Story 1

```bash
# Launch all unit tests for User Story 1 together (MANDATORY):
Task: "Unit test for ProfileRepository when saving profile persists to database in apps/backend/internal/repository/profile_repository_test.go"
Task: "Unit test for ReferenceLetterRepository when saving reference letter persists to database in apps/backend/internal/repository/reference_letter_repository_test.go"
Task: "Unit test for WorkExperienceRepository when saving work experience persists to database in apps/backend/internal/repository/work_experience_repository_test.go"
Task: "Unit test for CredibilityHighlightRepository when saving highlight persists to database in apps/backend/internal/repository/credibility_highlight_repository_test.go"

# Launch all repository interfaces together:
Task: "Create Profile repository interface in apps/backend/internal/repository/profile_repository.go"
Task: "Create ReferenceLetter repository interface in apps/backend/internal/repository/reference_letter_repository.go"
Task: "Create WorkExperience repository interface in apps/backend/internal/repository/work_experience_repository.go"
Task: "Create CredibilityHighlight repository interface in apps/backend/internal/repository/credibility_highlight_repository.go"

# Launch all repository implementations together:
Task: "Implement GormProfileRepository in apps/backend/internal/repository/gorm_profile_repository.go"
Task: "Implement GormReferenceLetterRepository in apps/backend/internal/repository/gorm_reference_letter_repository.go"
Task: "Implement GormWorkExperienceRepository in apps/backend/internal/repository/gorm_work_experience_repository.go"
Task: "Implement GormCredibilityHighlightRepository in apps/backend/internal/repository/gorm_credibility_highlight_repository.go"

# Launch frontend components together:
Task: "Create reference letter upload UI component in apps/frontend/components/profile/ReferenceLetterUpload.tsx"
Task: "Create profile generation trigger UI in apps/frontend/components/profile/GenerateProfileButton.tsx"
Task: "Create API client method for POST /reference-letters in apps/frontend/lib/api/referenceLetters.ts"
Task: "Create API client method for POST /profile/generate in apps/frontend/lib/api/profile.ts"
```

---

## Implementation Strategy

### MVP First (User Stories 1 & 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Profile Generation)
4. Complete Phase 4: User Story 2 (View Profile)
5. **STOP and VALIDATE**: Test User Stories 1 & 2 independently
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (Core MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo (Full Profile View)
4. Add User Story 3 → Test independently → Deploy/Demo (Tailoring)
5. Add User Story 4 → Test independently → Deploy/Demo (CV Download)
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Profile Generation)
   - Developer B: User Story 2 (View Profile) - can start after US1 models are done
   - Developer C: User Story 3 (Tailoring) - can start after US1/US2 are done
   - Developer D: User Story 4 (CV Download) - can start after US1/US2 are done
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **Unit tests are MANDATORY** - write them FIRST before implementation
- Follow AAA (Arrange-Act-Assert) style for all tests
- Use descriptive test names: `context_whenTrigger_thenExpectation`
- Stub external calls (HTTP, database) in tests - no real external dependencies
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Performance goal: Profile generation < 60s
- All file paths use absolute structure from monorepo root (apps/backend/, apps/frontend/)
- Always run `make test`, `make lint`, and `make fmt` before committing
