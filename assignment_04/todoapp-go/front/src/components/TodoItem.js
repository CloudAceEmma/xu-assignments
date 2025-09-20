import React, { useState } from 'react';

function TodoItem({ todo, onUpdateTodo, onDeleteTodo, onToggleTodoStatus }) {
  const [isEditing, setIsEditing] = useState(false);
  const [newTitle, setNewTitle] = useState(todo.title);
  const [newDescription, setNewDescription] = useState(todo.description || '');
  const [newPriority, setNewPriority] = useState(todo.priority || 0);
  const [newDueDate, setNewDueDate] = useState(
    todo.due_date ? new Date(todo.due_date).toISOString().slice(0, 16) : ''
  );

  const handleCheckboxChange = () => {
    onToggleTodoStatus(todo.id);
  };

  const handleTitleChange = (e) => {
    setNewTitle(e.target.value);
  };

  const handleEditClick = () => {
    setIsEditing(true);
  };

  const handleSaveClick = () => {
    if (newTitle.trim() !== '') {
      onUpdateTodo(todo.id, {
        title: newTitle,
        description: newDescription,
        priority: parseInt(newPriority, 10),
        due_date: newDueDate === '' ? null : new Date(newDueDate).toISOString(),
      });
      setIsEditing(false);
    }
  };

  const handleCancelClick = () => {
    setNewTitle(todo.title);
    setNewDescription(todo.description || '');
    setNewPriority(todo.priority || 0);
    setNewDueDate(
      todo.due_date ? new Date(todo.due_date).toISOString().slice(0, 16) : ''
    );
    setIsEditing(false);
  };

  const handleDeleteClick = () => {
    onDeleteTodo(todo.id);
  };

  return (
    <li className={`todo-item ${todo.completed ? 'completed' : ''}`}>
      <input
        type="checkbox"
        checked={todo.completed}
        onChange={handleCheckboxChange}
      />
      {isEditing ? (
        <div className="edit-mode-inputs">
          <input
            type="text"
            value={newTitle}
            onChange={handleTitleChange}
            onKeyPress={(e) => { if (e.key === 'Enter') handleSaveClick(); }}
          />
          <input
            type="text"
            value={newDescription}
            onChange={(e) => setNewDescription(e.target.value)}
            placeholder="Description"
          />
          <input
            type="number"
            value={newPriority}
            onChange={(e) => setNewPriority(e.target.value)}
            placeholder="Priority"
          />
          <input
            type="datetime-local"
            value={newDueDate}
            onChange={(e) => setNewDueDate(e.target.value)}
          />
          <button onClick={handleSaveClick}>Save</button>
          <button onClick={handleCancelClick}>Cancel</button>
        </div>
      ) : (
        <>
          <span onDoubleClick={handleEditClick}>
            {todo.title}
            {todo.description && ` - ${todo.description}`}
            {todo.priority > 0 && ` (P: ${todo.priority})`}
            {todo.due_date && ` (Due: ${new Date(todo.due_date).toLocaleString()})`}
          </span>
          <button onClick={handleEditClick}>Edit</button>
          <button onClick={handleDeleteClick}>Delete</button>
        </>
      )}
    </li>
  );
}

export default TodoItem;