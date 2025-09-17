from pydantic import BaseModel

class TodoBase(BaseModel):
    title: str

class TodoCreate(TodoBase):
    completed: bool = False

class TodoUpdate(TodoBase):
    completed: bool

class Todo(TodoBase):
    id: int
    completed: bool

    class Config:
        from_attributes = True
