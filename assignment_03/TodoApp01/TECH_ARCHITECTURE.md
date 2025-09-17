# Technical Architecture Document for To-Do Application

## 1. Technical Requirements

*   **Frontend:** React
*   **Backend:** FastAPI
*   **Database:** SQLite
*   **Project Structure:** The project will be organized into two main directories: `backend` and `front`.

## 2. SQL Table Design

The application will use a single table named `todos` to store to-do items.

```sql
CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT 0
);
```

## 3. API Interface Specifications

**Base URL:** `/api/v1`

### 3.1. Get All Todos

*   **Endpoint:** `GET /api/v1/todos`
*   **Description:** Retrieves all to-do items, with optional filtering by status.
*   **Query Parameters:**
    *   `status`: Optional. Filters to-do items by their completion status.
        *   `all`: Returns all to-do items (default).
        *   `incomplete`: Returns only incomplete to-do items.
        *   `completed`: Returns only completed to-do items.
*   **Response (200 OK):** A JSON array of to-do objects.
    ```json
    [
        {
            "id": 1,
            "title": "Learn FastAPI",
            "completed": false
        },
        {
            "id": 2,
            "title": "Build Todo App",
            "completed": true
        }
    ]
    ```

### 3.2. Create a New Todo

*   **Endpoint:** `POST /api/v1/todos`
*   **Description:** Adds a new to-do item to the database.
*   **Request Body (application/json):**
    ```json
    {
        "title": "New task title"
    }
    ```
*   **Response (201 Created):** The newly created to-do item.
    ```json
    {
        "id": 3,
        "title": "New task title",
        "completed": false
    }
    ```

### 3.3. Update a Todo

*   **Endpoint:** `PUT /api/v1/todos/{id}`
*   **Description:** Updates an existing to-do item identified by its ID. This can be used to mark a task as complete/incomplete or change its title.
*   **Path Parameters:**
    *   `id` (integer): The unique identifier of the to-do item to update.
*   **Request Body (application/json):**
    ```json
    {
        "title": "Updated task title",
        "completed": true
    }
    ```
*   **Response (200 OK):** The updated to-do item.
    ```json
    {
        "id": 1,
        "title": "Updated task title",
        "completed": true
    }
    ```
*   **Response (404 Not Found):** If the to-do item with the given ID does not exist.

### 3.4. Delete a Todo

*   **Endpoint:** `DELETE /api/v1/todos/{id}`
*   **Description:** Deletes a specific to-do item by its ID.
*   **Path Parameters:**
    *   `id` (integer): The unique identifier of the to-do item to delete.
*   **Response (200 OK):** A confirmation message.
    ```json
    {
        "message": "Todo item deleted successfully"
    }
    ```
*   **Response (404 Not Found):** If the to-do item with the given ID does not exist.

### 3.5. Clear Completed Todos

*   **Endpoint:** `DELETE /api/v1/todos/completed`
*   **Description:** Deletes all to-do items that are marked as completed.
*   **Response (200 OK):** A confirmation message.
    ```json
    {
        "message": "Completed todo items cleared successfully"
    }
    ```

### 3.6. Clear All Todos

*   **Endpoint:** `DELETE /api/v1/todos`
*   **Description:** Deletes all to-do items from the database.
*   **Response (200 OK):** A confirmation message.
    ```json
    {
        "message": "All todo items cleared successfully"
    }
    ```
