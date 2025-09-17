import React from 'react';

function FilterButtons({ currentFilter, onSetFilter }) {
  return (
    <div className="filter-buttons">
      <button
        className={currentFilter === 'all' ? 'active' : ''}
        onClick={() => onSetFilter('all')}
      >
        All
      </button>
      <button
        className={currentFilter === 'incomplete' ? 'active' : ''}
        onClick={() => onSetFilter('incomplete')}
      >
        Incomplete
      </button>
      <button
        className={currentFilter === 'completed' ? 'active' : ''}
        onClick={() => onSetFilter('completed')}
      >
        Completed
      </button>
    </div>
  );
}

export default FilterButtons;
