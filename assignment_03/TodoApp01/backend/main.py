from fastapi import FastAPI, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List, Optional
from backend import models, database
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="TodoApp API")

# CORS middleware to allow frontend to connect
origins = [
    "http://localhost",
    "http://localhost:3000", # Assuming React frontend runs on port 3000
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Create database tables on startup
@app.on_event("startup")
def on_startup():
    database.create_db_and_tables()

# Dependency to get DB session
def get_db():
    db = database.SessionLocal()
    try:
        yield db
    finally:
        db.close()

@app.get("/api/v1/todos", response_model=List[models.Todo])
def read_todos(status: Optional[str] = "all", db: Session = Depends(get_db)):
    if status == "completed":
        todos = db.query(database.Todo).filter(database.Todo.completed == True).all()
    elif status == "incomplete":
        todos = db.query(database.Todo).filter(database.Todo.completed == False).all()
    else:
        todos = db.query(database.Todo).all()
    return todos

@app.post("/api/v1/todos", response_model=models.Todo, status_code=201)
def create_todo(todo: models.TodoCreate, db: Session = Depends(get_db)):
    db_todo = database.Todo(title=todo.title, completed=todo.completed)
    db.add(db_todo)
    db.commit()
    db.refresh(db_todo)
    return db_todo

@app.put("/api/v1/todos/{todo_id}", response_model=models.Todo)
def update_todo(todo_id: int, todo: models.TodoUpdate, db: Session = Depends(get_db)):
    db_todo = db.query(database.Todo).filter(database.Todo.id == todo_id).first()
    if db_todo is None:
        raise HTTPException(status_code=404, detail="Todo not found")
    
    db_todo.title = todo.title
    db_todo.completed = todo.completed
    db.commit()
    db.refresh(db_todo)
    return db_todo

@app.delete("/api/v1/todos/{todo_id}")
def delete_todo(todo_id: int, db: Session = Depends(get_db)):
    db_todo = db.query(database.Todo).filter(database.Todo.id == todo_id).first()
    if db_todo is None:
        raise HTTPException(status_code=404, detail="Todo not found")
    db.delete(db_todo)
    db.commit()
    return {"message": "Todo item deleted successfully"}

@app.put("/api/v1/todos/clear-completed")
def clear_completed_todos(db: Session = Depends(get_db)):
    completed_todos = db.query(database.Todo).filter(database.Todo.completed == True).all()
    for todo in completed_todos:
        db.delete(todo)
    db.commit()
    return {"message": "Completed todo items cleared successfully"}

@app.delete("/api/v1/todos")
def clear_all_todos(db: Session = Depends(get_db)):
    db.query(database.Todo).delete()
    db.commit()
    return {"message": "All todo items cleared successfully"}