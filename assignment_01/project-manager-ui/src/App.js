import React from 'react';
import ProjectList from './components/ProjectList';
import './App.css'; // We will create this file next

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>项目任务管理系统</h1>
      </header>
      <nav className="App-sidebar">
        <ul>
          <li><a href="#">项目概览</a></li>
          <li><a href="#">我的任务</a></li>
          <li><a href="#">团队</a></li>
          <li><a href="#">设置</a></li>
        </ul>
      </nav>
      <main className="App-main-content">
        <ProjectList />
      </main>
    </div>
  );
}

export default App;