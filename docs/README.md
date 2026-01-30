# Amar Pathagar API Documentation

This directory contains the OpenAPI/Swagger documentation for the Amar Pathagar API.

## Files

- `swagger.yaml` - Complete OpenAPI 3.0 specification for all API endpoints

## Viewing the Documentation

### Option 1: Swagger UI (Recommended)

You can view the documentation using Swagger UI online:

1. Go to [Swagger Editor](https://editor.swagger.io/)
2. Click "File" → "Import file"
3. Select the `swagger.yaml` file
4. The documentation will be rendered with an interactive interface

### Option 2: Swagger UI Docker

Run Swagger UI locally using Docker:

```bash
docker run -p 8081:8080 -e SWAGGER_JSON=/docs/swagger.yaml -v $(pwd)/docs:/docs swaggerapi/swagger-ui
```

Then open http://localhost:8081 in your browser.

### Option 3: VS Code Extension

Install the "OpenAPI (Swagger) Editor" extension in VS Code and open the `swagger.yaml` file.

## API Overview

### Base URL

- Development: `http://localhost:8080/api/v1`
- Production: `https://api.amarpathagar.com/api/v1`

### Authentication

Most endpoints require authentication using JWT Bearer tokens.

To authenticate:
1. Register or login via `/auth/register` or `/auth/login`
2. Use the returned `access_token` in the Authorization header:
   ```
   Authorization: Bearer <your_access_token>
   ```

### Endpoint Categories

1. **Authentication** - User registration and login
2. **Users** - User profiles and management
3. **Books** - Book catalog and CRUD operations
4. **Book Requests** - Request and borrow books
5. **Handover** - Book delivery coordination
6. **Ideas** - Book suggestions and discussions
7. **Reviews** - Book reviews and ratings
8. **Bookmarks** - User bookmarks
9. **Donations** - Financial contributions
10. **Notifications** - User notifications
11. **Admin** - Administrative operations (admin role required)
12. **Leaderboard** - User rankings

## Key Features

### Trust-Based System

- Users have a **success score** that affects their ability to request books
- Minimum score of 20 required to request books
- Scores are earned through:
  - Returning books on time
  - Writing reviews
  - Contributing ideas
  - Community participation

### Book Handover System

1. User requests a book
2. Admin approves the request
3. Handover thread is created between current holder and requester
4. Users coordinate delivery through the thread
5. Receiver confirms delivery
6. Book status updates to "reading"

### Book Statuses

- `available` - Book is available for request
- `requested` - Book has been requested and approved, awaiting delivery
- `reading` - Book is currently being read
- `on_hold` - Book reading completed, awaiting next request
- `reserved` - Book is reserved (future use)

## Testing the API

### Using cURL

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }'

# Get books (with auth token)
curl -X GET http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer <your_token>"
```

### Using Postman

1. Import the `swagger.yaml` file into Postman
2. Postman will automatically create a collection with all endpoints
3. Set up an environment variable for the auth token
4. Start testing!

## Response Format

### Success Response

```json
{
  "data": {
    // Response data here
  }
}
```

### Error Response

```json
{
  "error": "Error message here"
}
```

## Rate Limiting

Currently, there are no rate limits implemented. This may change in production.

## Support

For issues or questions about the API, please contact the development team or open an issue on GitHub.
