import React, { useState, useEffect } from 'react';
import TodoList from './components/TodoList';
import TodoForm from './components/TodoForm';
import FilterButtons from './components/FilterButtons';
import { getTodos, createTodo, updateTodo, deleteTodo, clearCompletedTodos, clearAllTodos, toggleTodoStatus } from './api';
import './App.css'; // Assuming you'll create an App.css for styling

function App() {
  const [todos, setTodos] = useState([]);
  const [filter, setFilter] = useState('all'); // 'all', 'incomplete', 'completed'

  useEffect(() => {
    fetchTodos();
  }, [filter]);

  const fetchTodos = async () => {
    try {
      const data = await getTodos(filter);
      setTodos(data);
    } catch (error) {
      console.error("Error fetching todos:", error);
    }
  };

  const handleAddTodo = async (todo) => {
    try {
      const newTodo = await createTodo(todo);
      setTodos([...todos, newTodo]);
    } catch (error) {
      console.error("Error adding todo:", error);
    }
  };

  const handleUpdateTodo = async (id, updatedFields) => {
    try {
      const updatedTodo = await updateTodo(id, updatedFields);
      setTodos(todos.map(todo => (todo.id === id ? updatedTodo : todo)));
    } catch (error) {
      console.error("Error updating todo:", error);
    }
  };

  const handleDeleteTodo = async (id) => {
    try {
      await deleteTodo(id);
      setTodos(todos.filter(todo => todo.id !== id));
    } catch (error) {
      console.error("Error deleting todo:", error);
    }
  };

  const handleToggleTodoStatus = async (id) => {
    try {
      await toggleTodoStatus(id);
      fetchTodos(); // Re-fetch todos to get the updated list
    } catch (error) {
      console.error("Error toggling todo status:", error);
    }
  };

  const handleClearCompleted = async () => {
    try {
      await clearCompletedTodos();
      fetchTodos(); // Re-fetch todos to get the updated list
    } catch (error) {
      console.error("Error clearing completed todos:", error);
    }
  };

  const handleClearAll = async () => {
    try {
      await clearAllTodos();
      setTodos([]); // Clear all todos from the state
    } catch (error) {
      console.error("Error clearing all todos:", error);
    }
  };

  return (
    <div className="App">
      <h1>Todo Application</h1>
      <TodoForm onAddTodo={handleAddTodo} />
      <FilterButtons currentFilter={filter} onSetFilter={setFilter} />
      <TodoList
        todos={todos}
        onUpdateTodo={handleUpdateTodo}
        onDeleteTodo={handleDeleteTodo}
        onToggleTodoStatus={handleToggleTodoStatus}
      />
      <div className="actions">
        <button onClick={handleClearCompleted}>Clear Completed</button>
        <button onClick={handleClearAll}>Clear All</button>
      </div>
    </div>
  );
}

export default App;
