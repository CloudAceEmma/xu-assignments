import React, { useState } from 'react';

function TodoItem({ todo, onUpdateTodo, onDeleteTodo }) {
  const [isEditing, setIsEditing] = useState(false);
  const [newTitle, setNewTitle] = useState(todo.title);

  const handleCheckboxChange = () => {
    onUpdateTodo(todo.id, { completed: !todo.completed, title: todo.title });
  };

  const handleTitleChange = (e) => {
    setNewTitle(e.target.value);
  };

  const handleEditClick = () => {
    setIsEditing(true);
  };

  const handleSaveClick = () => {
    if (newTitle.trim() !== '') {
      onUpdateTodo(todo.id, { title: newTitle, completed: todo.completed });
      setIsEditing(false);
    }
  };

  const handleCancelClick = () => {
    setNewTitle(todo.title);
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
        <>
          <input
            type="text"
            value={newTitle}
            onChange={handleTitleChange}
            onKeyPress={(e) => { if (e.key === 'Enter') handleSaveClick(); }}
          />
          <button onClick={handleSaveClick}>Save</button>
          <button onClick={handleCancelClick}>Cancel</button>
        </>
      ) : (
        <>
          <span onDoubleClick={handleEditClick}>{todo.title}</span>
          <button onClick={handleEditClick}>Edit</button>
          <button onClick={handleDeleteClick}>Delete</button>
        </>
      )}
    </li>
  );
}

export default TodoItem;
