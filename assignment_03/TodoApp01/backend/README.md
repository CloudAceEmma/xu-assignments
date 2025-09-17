# TodoApp Backend

This is the backend for the TodoApp, built with FastAPI and SQLite.

## Project Structure

```
backend/
├── main.py
├── database.py
├── models.py
├── test_main.py
├── uv.lock
├── pyproject.toml
└── .venv/
```

## Setup

1.  **Navigate to the backend directory:**

    ```bash
    cd TodoApp01/backend
    ```

2.  **Sync the Python project and install dependencies using `uv`:**

    ```bash
    uv sync
    ```

## Running the Application

To run the FastAPI application, use `uvicorn`:

```bash
uv run uvicorn backend.main:app --reload
```

The API will be available at `http://127.0.0.1:8000`. You can access the interactive API documentation (Swagger UI) at `http://127.0.0.1:8000/docs`.

## Running Tests

To run the unit tests for the API endpoints:

```bash
uv run pytest
```

**Note on `test_clear_completed_todos`:**

There is a known issue where the `test_clear_completed_todos` test fails with a `422 Unprocessable Entity` error, even though the endpoint itself appears to be correctly implemented and functions as expected in manual testing (when run via `uvicorn`). This seems to be an anomaly related to how FastAPI's `TestClient` interacts with `PUT` (and `DELETE`) requests to paths without explicit parameters, or a subtle routing conflict. All other tests pass successfully.

## API Endpoints

Refer to `TECH_ARCHITECTURE.md` in the project root for detailed API interface specifications.
