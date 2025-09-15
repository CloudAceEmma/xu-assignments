import React from 'react';
import './ProjectCard.css';

const ProjectCard = ({ project }) => {
  return (
    <div className="project-card">
      <h3 className="project-card-title">{project.name}</h3>
      <p className="project-card-detail">部门: {project.department}</p>
      <p className="project-card-detail">状态: <span className={`status-${project.status.toLowerCase()}`}>{project.status}</span></p>
      <p className="project-card-detail">完成率: {project.completionRate}%</p>
      <div className="project-card-progress-bar">
        <div className="project-card-progress-fill" style={{ width: `${project.completionRate}%` }}></div>
      </div>
    </div>
  );
};

export default ProjectCard;
