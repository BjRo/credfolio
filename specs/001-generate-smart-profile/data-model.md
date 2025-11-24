# Data Model: Generate Smart Profile

## Entities

### User (Existing/Mock)
- **ID**: UUID (PK)
- **Email**: String
- **Name**: String

### Profile
- **ID**: UUID (PK)
- **UserID**: UUID (FK -> User.ID)
- **Summary**: Text
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp

### WorkExperience
- **ID**: UUID (PK)
- **ProfileID**: UUID (FK -> Profile.ID)
- **CompanyName**: String
- **Role**: String
- **StartDate**: Date
- **EndDate**: Date (Nullable)
- **Description**: Text
- **ReferenceLetterID**: UUID (FK -> ReferenceLetter.ID, Nullable) - *Link to source*

### Skill
- **ID**: UUID (PK)
- **Name**: String (Unique)

### ProfileSkill
- **ProfileID**: UUID (FK)
- **SkillID**: UUID (FK)
- **Proficiency**: String (Optional)

### ReferenceLetter
- **ID**: UUID (PK)
- **UserID**: UUID (FK -> User.ID)
- **FileName**: String
- **StoragePath**: String
- **UploadDate**: Timestamp
- **Status**: String (PENDING, PROCESSED, FAILED)
- **ExtractedText**: Text (for re-processing)

### CredibilityHighlight
- **ID**: UUID (PK)
- **WorkExperienceID**: UUID (FK -> WorkExperience.ID)
- **Quote**: Text
- **Sentiment**: String (POSITIVE, NEUTRAL)
- **SourceLetterID**: UUID (FK -> ReferenceLetter.ID)

### JobMatch (Tailored Profile)
- **ID**: UUID (PK)
- **BaseProfileID**: UUID (FK -> Profile.ID)
- **JobDescription**: Text
- **MatchScore**: Float
- **TailoredSummary**: Text
- **CreatedAt**: Timestamp

## Relationships

- User **HasOne** Profile
- Profile **HasMany** WorkExperience
- Profile **ManyToMany** Skill
- User **HasMany** ReferenceLetter
- WorkExperience **HasMany** CredibilityHighlight
- Profile **HasMany** JobMatch

## GORM Struct Preview

```go
type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Profile   Profile
}

type Profile struct {
    ID              uuid.UUID `gorm:"type:uuid;primary_key"`
    UserID          uuid.UUID
    WorkExperiences []WorkExperience
    Skills          []*Skill `gorm:"many2many:profile_skills;"`
}

type WorkExperience struct {
    ID                  uuid.UUID `gorm:"type:uuid;primary_key"`
    ProfileID           uuid.UUID
    CredibilityHighlights []CredibilityHighlight
}
```

