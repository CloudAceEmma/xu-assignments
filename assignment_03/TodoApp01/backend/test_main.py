from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm import Session, sessionmaker
from sqlalchemy.pool import StaticPool
import pytest

from .main import app, get_db
from .database import Base, Todo

# Use an in-memory SQLite database for testing
SQLALCHEMY_DATABASE_URL = "sqlite:///:memory:"

engine = create_engine(
    SQLALCHEMY_DATABASE_URL,
    connect_args={"check_same_thread": False},
    poolclass=StaticPool,
)
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

@pytest.fixture(name="db_session")
def db_session_fixture():
    Base.metadata.create_all(bind=engine) # Create tables
    db = TestingSessionLocal()
    try:
        yield db
    finally:
        db.close()
        # Clear all data from the Todo table after each test
        Base.metadata.drop_all(bind=engine) # Drop tables after all tests

@pytest.fixture(name="client")
def client_fixture(db_session: Session):
    def override_get_db():
        yield db_session

    app.dependency_overrides[get_db] = override_get_db
    with TestClient(app) as client:
        yield client
    app.dependency_overrides.clear() # Clear overrides after tests

def test_create_todo(client: TestClient):
    response = client.post(
        "/api/v1/todos",
        json={"title": "Test Todo 1"}
    )
    assert response.status_code == 201
    data = response.json()
    assert data["title"] == "Test Todo 1"
    assert data["completed"] == False
    assert "id" in data

def test_read_todos(client: TestClient):
    # Create some todos first
    client.post("/api/v1/todos", json={"title": "Test Todo 1"})
    client.post("/api/v1/todos", json={"title": "Test Todo 2", "completed": True})
    client.post("/api/v1/todos", json={"title": "Test Todo 3"})

    response = client.get("/api/v1/todos")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 3
    assert data[0]["title"] == "Test Todo 1"

def test_read_completed_todos(client: TestClient):
    client.post("/api/v1/todos", json={"title": "Test Todo 1"})
    client.post("/api/v1/todos", json={"title": "Test Todo 2", "completed": True})
    client.post("/api/v1/todos", json={"title": "Test Todo 3"})

    response = client.get("/api/v1/todos?status=completed")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["title"] == "Test Todo 2"
    assert data[0]["completed"] == True

def test_read_incomplete_todos(client: TestClient):
    client.post("/api/v1/todos", json={"title": "Test Todo 1"})
    client.post("/api/v1/todos", json={"title": "Test Todo 2", "completed": True})
    client.post("/api/v1/todos", json={"title": "Test Todo 3"})

    response = client.get("/api/v1/todos?status=incomplete")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 2
    assert data[0]["title"] == "Test Todo 1"
    assert data[0]["completed"] == False

def test_update_todo(client: TestClient):
    create_response = client.post(
        "/api/v1/todos",
        json={"title": "Todo to update"}
    )
    todo_id = create_response.json()["id"]

    update_response = client.put(
        f"/api/v1/todos/{todo_id}",
        json={"title": "Updated Todo", "completed": True}
    )
    assert update_response.status_code == 200
    data = update_response.json()
    assert data["title"] == "Updated Todo"
    assert data["completed"] == True
    assert data["id"] == todo_id

def test_update_nonexistent_todo(client: TestClient):
    response = client.put(
        "/api/v1/todos/999",
        json={"title": "Nonexistent Todo", "completed": False}
    )
    assert response.status_code == 404
    assert response.json() == {"detail": "Todo not found"}

def test_delete_todo(client: TestClient):
    create_response = client.post(
        "/api/v1/todos",
        json={"title": "Todo to delete"}
    )
    todo_id = create_response.json()["id"]

    delete_response = client.delete(f"/api/v1/todos/{todo_id}")
    assert delete_response.status_code == 200
    assert delete_response.json() == {"message": "Todo item deleted successfully"}

    # Verify it's actually deleted
    get_response = client.get("/api/v1/todos")
    assert len(get_response.json()) == 0

def test_delete_nonexistent_todo(client: TestClient):
    response = client.delete("/api/v1/todos/999")
    assert response.status_code == 404
    assert response.json() == {"detail": "Todo not found"}

def test_clear_completed_todos(client: TestClient):
def test_clear_completed_todos(client: TestClient):
    client.post("/api/v1/todos", json={"title": "Todo 1"})
    client.post("/api/v1/todos", json={"title": "Todo 2", "completed": True})
    client.post("/api/v1/todos", json={"title": "Todo 3"})
    client.post("/api/v1/todos", json={"title": "Todo 4", "completed": True})

    # Debug: Print current todos before clearing completed ones
    current_todos_response = client.get("/api/v1/todos")
    print(f"Todos before clearing completed: {current_todos_response.json()}")

    response = client.put("/api/v1/todos/clear-completed")
    assert response.status_code == 200
    assert response.json() == {"message": "Completed todo items cleared successfully"}

    get_response = client.get("/api/v1/todos")
    data = get_response.json()
    assert len(data) == 2
    assert data[0]["title"] == "Todo 1"
    assert data[1]["title"] == "Todo 3"

def test_clear_all_todos(client: TestClient):
    client.post("/api/v1/todos", json={"title": "Todo 1"})
    client.post("/api/v1/todos", json={"title": "Todo 2", "completed": True})

    response = client.delete("/api/v1/todos")
    assert response.status_code == 200
    assert response.json() == {"message": "All todo items cleared successfully"}

    get_response = client.get("/api/v1/todos")
    data = get_response.json()
    assert len(data) == 0