import React from 'react';
import ProjectCard from './ProjectCard';
import './ProjectList.css';

const ProjectList = () => {
  const projects = [
    {
      id: 1,
      name: 'AI模型优化项目',
      department: '研发部',
      status: '进行中',
      completionRate: 75,
    },
    {
      id: 2,
      name: '新产品发布计划',
      department: '市场部',
      status: '未开始',
      completionRate: 0,
    },
    {
      id: 3,
      name: '内部系统升级',
      department: 'IT部',
      status: '已完成',
      completionRate: 100,
    },
    {
      id: 4,
      name: '用户体验改进',
      department: '产品部',
      status: '进行中',
      completionRate: 40,
    },
    {
      id: 5,
      name: '季度财报分析',
      department: '财务部',
      status: '进行中',
      completionRate: 90,
    },
  ];

  return (
    <div className="project-list-container">
      <h2 className="project-list-title">项目概览</h2>
      <div className="project-list-grid">
        {projects.map((project) => (
          <ProjectCard key={project.id} project={project} />
        ))}
      </div>
    </div>
  );
};

export default ProjectList;
