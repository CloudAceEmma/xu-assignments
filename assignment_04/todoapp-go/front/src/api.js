const API_BASE_URL = 'http://localhost:8000/api/v1';

export const getTodos = async (filter = 'all') => {
  let url = `${API_BASE_URL}/todos`;
  if (filter !== 'all') {
    const completed = filter === 'completed';
    url += `?completed=${completed}`;
  }
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error('Failed to fetch todos');
  }
  const data = await response.json();
  return data.data; // Extract the actual todo array from the response
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
  const data = await response.json();
  return data.data; // Extract the created todo object from the response
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
  const data = await response.json();
  return data.data; // Extract the updated todo object from the response
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

export const toggleTodoStatus = async (id) => {
  const response = await fetch(`${API_BASE_URL}/todos/${id}/toggle`, {
    method: 'PATCH',
  });
  if (!response.ok) {
    throw new Error('Failed to toggle todo status');
  }
  const data = await response.json();
  return data.data; // The backend returns the updated todo data
};

export const clearCompletedTodos = async () => {
  const response = await fetch(`${API_BASE_URL}/todos/completed`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to clear completed todos');
  }
  return response.json();
};

export const clearAllTodos = async () => {
  const response = await fetch(`${API_BASE_URL}/todos/all`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('Failed to clear all todos');
  }
  return response.json();
};
