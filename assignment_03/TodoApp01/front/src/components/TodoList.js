import React from 'react';
import TodoItem from './TodoItem';

function TodoList({ todos, onUpdateTodo, onDeleteTodo }) {
  return (
    <div className="todo-list">
      {todos.length === 0 ? (
        <p>No todos found.</p>
      ) : (
        <ul>
          {todos.map(todo => (
            <TodoItem
              key={todo.id}
              todo={todo}
              onUpdateTodo={onUpdateTodo}
              onDeleteTodo={onDeleteTodo}
            />
          ))}
        </ul>
      )}
    </div>
  );
}

export default TodoList;
