# credfolio

## API Documentation

### Upload Reference Letter
- **POST** `/api/upload`
- **Body**: Multipart form data, key `file` (PDF)
- **Response**: `{ "status": "processing", "user_id": "uuid" }`

### Get Profile
- **GET** `/api/profile?user_id={uuid}`
- **Response**: User Profile JSON

### Download CV
- **GET** `/api/profile/cv?user_id={uuid}`
- **Response**: PDF File

### Tailor Profile
- **POST** `/api/profile/tailor`
- **Body**: `{ "user_id": "uuid", "job_description": "text" }`
- **Response**: Tailoring Result JSON
