---
description: "Task list for Generate Smart Profile & Credibility feature implementation"
---

# Tasks: Generate Smart Profile & Credibility

**Input**: Design documents from `/specs/001-generate-smart-profile/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Follow the ADRs, add tests for every task. Make sure to run `make fmt` and that `make lint` and `make test` pass for every task.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Monorepo**: `apps/backend/`, `apps/frontend/`
- **Backend**: `apps/backend/internal/{domain,service,repository,handler}`, `apps/backend/pkg/{ai,pdf,extraction}`
- **Frontend**: `apps/frontend/app/profile`, `apps/frontend/components/profile`, `apps/frontend/lib/api`

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

- [ ] T008 Setup PostgreSQL database connection using GORM in apps/backend/internal/repository/db.go
- [ ] T009 Create base User model (mock/existing) in apps/backend/internal/domain/user.go
- [ ] T010 [P] Create base Profile model in apps/backend/internal/domain/profile.go
- [ ] T011 [P] Create base WorkExperience model in apps/backend/internal/domain/work_experience.go
- [ ] T012 [P] Create base Skill model in apps/backend/internal/domain/skill.go
- [ ] T013 [P] Create base ReferenceLetter model in apps/backend/internal/domain/reference_letter.go
- [ ] T014 [P] Create base CredibilityHighlight model in apps/backend/internal/domain/credibility_highlight.go
- [ ] T015 [P] Create base JobMatch model in apps/backend/internal/domain/job_match.go
- [ ] T016 Setup GORM AutoMigrate for all domain models in apps/backend/internal/repository/migrations.go
- [ ] T017 [P] Create LLMProvider interface abstraction in apps/backend/internal/service/llm_provider.go
- [ ] T018 [P] Implement OpenAIProvider struct satisfying LLMProvider interface in apps/backend/pkg/ai/openai_provider.go
- [ ] T019 [P] Create text extraction service using in apps/backend/pkg/extraction/extractor.go
- [ ] T020 [P] Create PDF generation service using maroto in apps/backend/pkg/pdf/generator.go
- [ ] T021 Setup Chi router and middleware structure in apps/backend/cmd/server/main.go
- [ ] T022 Create mock authentication middleware (inject current user ID) in apps/backend/internal/handler/middleware/auth.go
- [ ] T023 Configure error handling and logging infrastructure in apps/backend/pkg/logger/logger.go
- [ ] T024 Generate OpenAPI server stubs using oapi-codegen in apps/backend/api/generated/
- [ ] T025 Generate TypeScript API client from OpenAPI spec in apps/frontend/lib/api/generated/

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Profile Generation from References (Priority: P1) 🎯 MVP

**Goal**: Automatically populate professional profile by extracting data from uploaded reference letters using AI, including credibility highlights from employer sentiment.

**Independent Test**: Upload a sample reference letter and verify that the "Experience" and "Skills" sections are populated with correct data, and that "Credibility" section contains positive quotes or sentiment summaries.

### Implementation for User Story 1

- [ ] T026 [P] [US1] Create Profile repository interface in apps/backend/internal/repository/profile_repository.go
- [ ] T027 [P] [US1] Create ReferenceLetter repository interface in apps/backend/internal/repository/reference_letter_repository.go
- [ ] T028 [P] [US1] Implement GormProfileRepository in apps/backend/internal/repository/gorm_profile_repository.go
- [ ] T029 [P] [US1] Implement GormReferenceLetterRepository in apps/backend/internal/repository/gorm_reference_letter_repository.go
- [ ] T030 [US1] Create ProfileService with GenerateProfileFromReferences method in apps/backend/internal/service/profile_service.go
- [ ] T031 [US1] Implement AI extraction logic in ProfileService that calls LLMProvider to extract structured data (Company, Role, Dates, Skills, Achievements) in apps/backend/internal/service/profile_service.go
- [ ] T032 [US1] Implement credibility extraction logic that extracts positive sentiment quotes from reference letters in apps/backend/internal/service/profile_service.go
- [ ] T033 [US1] Create WorkExperience repository interface in apps/backend/internal/repository/work_experience_repository.go
- [ ] T034 [US1] Implement GormWorkExperienceRepository in apps/backend/internal/repository/gorm_work_experience_repository.go
- [ ] T035 [US1] Create CredibilityHighlight repository interface in apps/backend/internal/repository/credibility_highlight_repository.go
- [ ] T036 [US1] Implement GormCredibilityHighlightRepository in apps/backend/internal/repository/gorm_credibility_highlight_repository.go
- [ ] T037 [US1] Implement file upload handler for reference letters (multipart/form-data) in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T038 [US1] Implement POST /reference-letters endpoint handler in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T039 [US1] Implement POST /profile/generate endpoint handler that processes uploaded reference letters in apps/backend/internal/handler/profile_handler.go
- [ ] T040 [US1] Add validation and error handling for file uploads and AI extraction in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T041 [US1] Add logging for profile generation operations in apps/backend/internal/service/profile_service.go
- [ ] T042 [US1] Create reference letter upload UI component in apps/frontend/components/profile/ReferenceLetterUpload.tsx
- [ ] T043 [US1] Create profile generation trigger UI in apps/frontend/components/profile/GenerateProfileButton.tsx
- [ ] T044 [US1] Create API client method for POST /reference-letters in apps/frontend/lib/api/referenceLetters.ts
- [ ] T045 [US1] Create API client method for POST /profile/generate in apps/frontend/lib/api/profile.ts
- [ ] T046 [US1] Create profile generation page in apps/frontend/app/profile/generate/page.tsx
- [ ] T047 [US1] Implement profile data editing interface (edit/delete extracted information) in apps/frontend/components/profile/ProfileEditor.tsx

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - users can upload reference letters, generate profiles, and edit extracted data.

---

## Phase 4: User Story 2 - View Credibility Profile (Priority: P1)

**Goal**: Display profile in a comprehensive "LinkedIn-on-steroids" format that highlights employer feedback and credibility highlights, with aggregated skills and endorsements across experiences.

**Independent Test**: Navigate to the profile view after data population and verify standard sections (Experience, Skills) are augmented with "Employer Feedback" or "Credibility Highlights", and that skills/endorsements are aggregated across multiple reference letters.

### Implementation for User Story 2

- [ ] T048 [US2] Implement GET /profile endpoint handler in apps/backend/internal/handler/profile_handler.go
- [ ] T049 [US2] Add logic to aggregate skills and endorsements across multiple work experiences in ProfileService.GetProfile method in apps/backend/internal/service/profile_service.go
- [ ] T050 [US2] Create API client method for GET /profile in apps/frontend/lib/api/profile.ts
- [ ] T051 [US2] Create ProfileView component displaying Experience, Skills, and Credibility Highlights in apps/frontend/components/profile/ProfileView.tsx
- [ ] T052 [US2] Create CredibilityHighlights section component in apps/frontend/components/profile/CredibilityHighlights.tsx
- [ ] T053 [US2] Create WorkExperience display component with credibility highlights in apps/frontend/components/profile/WorkExperienceCard.tsx
- [ ] T054 [US2] Create Skills aggregation display component in apps/frontend/components/profile/SkillsSection.tsx
- [ ] T055 [US2] Create profile view page in apps/frontend/app/profile/page.tsx
- [ ] T056 [US2] Add styling for "LinkedIn-on-steroids" profile layout using Tailwind CSS in apps/frontend/components/profile/ProfileView.tsx

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - users can generate profiles and view them with credibility highlights.

---

## Phase 5: User Story 3 - Tailor Profile to Job Description (Priority: P2)

**Goal**: Generate a custom version of profile and CV that emphasizes the most relevant experience and skills based on semantic similarity to a provided job description, with match score and explanation.

**Independent Test**: Upload a Job Description and compare the standard profile vs. the tailored profile, verifying that skills/experiences matching JD keywords are reordered or highlighted, and that a Match Score or explanation is displayed.

### Implementation for User Story 3

- [ ] T057 [US3] Create JobMatch repository interface in apps/backend/internal/repository/job_match_repository.go
- [ ] T058 [US3] Implement GormJobMatchRepository in apps/backend/internal/repository/gorm_job_match_repository.go
- [ ] T059 [US3] Create TailoringService with TailorProfileToJobDescription method in apps/backend/internal/service/tailoring_service.go
- [ ] T060 [US3] Implement semantic matching logic using LLMProvider to rank experience/skills based on job description in apps/backend/internal/service/tailoring_service.go
- [ ] T061 [US3] Implement match score calculation in TailoringService in apps/frontend/lib/api/profile.ts
- [ ] T062 [US3] Implement POST /profile/tailor endpoint handler in apps/backend/internal/handler/profile_handler.go
- [ ] T063 [US3] Add validation for job description input in apps/backend/internal/handler/profile_handler.go
- [ ] T064 [US3] Create API client method for POST /profile/tailor in apps/frontend/lib/api/profile.ts
- [ ] T065 [US3] Create JobDescriptionInput component in apps/frontend/components/profile/JobDescriptionInput.tsx
- [ ] T066 [US3] Create TailoredProfileView component showing highlighted/reordered content in apps/frontend/components/profile/TailoredProfileView.tsx
- [ ] T067 [US3] Create MatchScore display component in apps/frontend/components/profile/MatchScore.tsx
- [ ] T068 [US3] Create profile tailoring page in apps/frontend/app/profile/tailor/page.tsx
- [ ] T069 [US3] Add explanation UI for why certain elements are highlighted in apps/frontend/components/profile/TailoredProfileView.tsx

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently - users can generate profiles, view them, and tailor them to job descriptions.

---

## Phase 6: User Story 4 - Download CV (Priority: P2)

**Goal**: Generate and download a formatted PDF CV from the current profile view (Standard or Tailored), reflecting emphasized content when tailored.

**Independent Test**: Generate a PDF file from a standard or tailored profile and verify the PDF contains all profile data in a formatted CV layout, with tailored content emphasized when applicable.

### Implementation for User Story 4

- [ ] T070 [US4] Enhance PDF generation service to create CV layout using maroto in apps/backend/pkg/pdf/cv_generator.go
- [ ] T071 [US4] Implement CV template with sections: Summary, Work Experience, Skills, Credibility Highlights in apps/backend/pkg/pdf/cv_generator.go
- [ ] T072 [US4] Add logic to emphasize tailored content in PDF when JobMatch is provided in apps/backend/pkg/pdf/cv_generator.go
- [ ] T073 [US4] Implement GET /profile/{profileId}/cv endpoint handler in apps/backend/internal/handler/profile_handler.go
- [ ] T074 [US4] Add query parameter support for tailored CV (jobMatchId) in GET /profile/{profileId}/cv handler in apps/backend/internal/handler/profile_handler.go
- [ ] T075 [US4] Create API client method for GET /profile/{profileId}/cv in apps/frontend/lib/api/profile.ts
- [ ] T076 [US4] Create DownloadCVButton component in apps/frontend/components/profile/DownloadCVButton.tsx
- [ ] T077 [US4] Add download CV functionality to profile view page in apps/frontend/app/profile/page.tsx
- [ ] T078 [US4] Add download CV functionality to tailored profile view in apps/frontend/app/profile/tailor/page.tsx
- [ ] T079 [US4] Add error handling for PDF generation failures in apps/backend/internal/handler/profile_handler.go

**Checkpoint**: All user stories should now be independently functional - users can generate profiles, view them with credibility, tailor them to jobs, and download CVs.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T080 [P] Update API documentation in apps/backend/api/openapi.yaml with complete request/response examples
- [ ] T081 [P] Add comprehensive error messages and user feedback throughout frontend components
- [ ] T082 Code cleanup and refactoring: Extract common patterns in service layer
- [ ] T083 Performance optimization: Cache AI responses where appropriate in apps/backend/internal/service/profile_service.go
- [ ] T084 [P] Add loading states and progress indicators for AI operations in apps/frontend/components/profile/
- [ ] T085 [P] Add input validation and sanitization for all user inputs in apps/backend/internal/handler/
- [ ] T086 Security hardening: Validate file types and sizes for reference letter uploads in apps/backend/internal/handler/reference_letter_handler.go
- [ ] T087 Run quickstart.md validation: Verify all setup steps work end-to-end
- [ ] T088 [P] Add comprehensive logging for all API endpoints in apps/backend/internal/handler/
- [ ] T089 Optimize database queries with proper indexing in apps/backend/internal/repository/
- [ ] T090 Add rate limiting for AI API calls in apps/backend/internal/service/profile_service.go

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
- Models within a story marked [P] can run in parallel
- Repository implementations marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members (after foundational phase)

---

## Parallel Example: User Story 1

```bash
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
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Performance goal: Profile generation < 60s
- All file paths use absolute structure from monorepo root (apps/backend/, apps/frontend/)

