const API_BASE_URL = 'http://localhost:8000/api/v1';

export const getTodos = async (status = 'all') => {
  const response = await fetch(`${API_BASE_URL}/todos?status=${status}`);
  if (!response.ok) {
    throw new Error('Failed to fetch todos');
  }
  return response.json();
};

export const createTodo = async (todo) => {
  const response = await fetch(`${API_BASE_URL}/todos`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(todo),
  });
  if (!response.ok) {
    throw new Error('Failed to create todo');
  }
  return response.json();
};

export const updateTodo = async (id, updatedFields) => {
  const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(updatedFields),
  });
  if (!response.ok) {
    throw new Error('Failed to update todo');
  }
  return response.json();
};

export const deleteTodo = async (id) => {
  const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to delete todo');
  }
  return response.json();
};

export const clearCompletedTodos = async () => {
  const response = await fetch(`${API_BASE_URL}/todos/clear-completed`, {
    method: 'PUT',
  });
  if (!response.ok) {
    throw new Error('Failed to clear completed todos');
  }
  return response.json();
};

export const clearAllTodos = async () => {
  const response = await fetch(`${API_BASE_URL}/todos`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to clear all todos');
  }
  return response.json();
};
