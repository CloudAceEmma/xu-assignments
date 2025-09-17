Create a To-Do Application
1. It needs to include a title, a form for entering new tasks (including an input field and an add button), and an ordered list for displaying the to-do items.
2. Please design a modern, concise CSS style. The main body should be centered with a maximum width of 800px. The input field and button styles should be aesthetically pleasing. There should be spacing between list items. When hovering over buttons and list items, there should be visual feedback.
3. Implement the add to-do item function. When the user clicks the add button, take the content from the input field, create a new `li` and add it to the list. After adding, clear the input field.
4. Implement mark as complete and delete functions. Inside each list item (`li`), there should be a 'Complete' button and a 'Delete' button. Clicking the 'Complete' button should add a 'completed' CSS class to the list item (to strike through the text). Clicking the 'Delete' button should remove the list item from the list.
5. Implement filtering conditions: "All", "Incomplete", "Completed". Implement bottom functions: "Clear Completed", "Clear All".

Technical Requirements: Frontend: React, Backend: FastAPI, Database: SQLite
Project Structure: `backend`, `front` 2 directories
Table Design: Please design a table for the to-do application and provide the SQL statement.