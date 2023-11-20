Certainly! Based on the provided information, here's a sample README file with Swagger API documentation:

```markdown
# Thoth Backend

Thoth Backend is a sample application built with the Echo framework in Go.

## API Documentation

API documentation is generated using Swagger. You can view it by visiting [Swagger Documentation](./docs/index.html).

### Create a new blog

- **URL:** `/blog/create`
- **Method:** `POST`
- **Description:** Create a new blog with the provided data
- **Request Body:**
  ```json
  {
    "title": "Sample Blog",
    "text": "This is a sample blog post.",
    "categories": ["Technology", "Programming"]
  }
  ```
- **Response:**
  - Status Code: 201
  - Body: `"Created"`

### Publish a blog

- **URL:** `/blog/publish/:id`
- **Method:** `POST`
- **Description:** Publish a blog with the specified ID
- **Request Parameters:**
  - `id`: Blog ID
- **Response:**
  - Status Code: 200
  - Body: (Empty)

### Add a comment to a blog

- **URL:** `/blog/comments/:id`
- **Method:** `POST`
- **Description:** Add a comment to the blog with the specified ID
- **Request Body:**
  ```json
  {
    "commenterName": "John Doe",
    "text": "Great blog post!"
  }
  ```
- **Response:**
  - Status Code: 200
  - Body: (Empty)

### Get comments for a blog

- **URL:** `/blog/comments/:id`
- **Method:** `GET`
- **Description:** Get comments for the blog with the specified ID
- **Query Parameters:**
  - `page`: Page number (optional, default is 1)
- **Response:**
  - Status Code: 200
  - Body: List of comments in JSON format

### Approve a comment

- **URL:** `/blog/comments/:id/approve`
- **Method:** `POST`
- **Description:** Approve a comment with the specified ID
- **Request Parameters:**
  - `id`: Comment ID
- **Response:**
  - Status Code: 200
  - Body: (Empty)

### Delete a comment

- **URL:** `/blog/comments/:id`
- **Method:** `POST`
- **Description:** Delete a comment with the specified ID
- **Request Parameters:**
  - `id`: Comment ID
- **Response:**
  - Status Code: 200
  - Body: (Empty)


