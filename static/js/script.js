document.addEventListener('DOMContentLoaded', () => {
  const addTaskBtn = document.getElementById('addTaskBtn');
  const taskInput = document.getElementById('taskInput');
  const taskList = document.getElementById('taskList');

  fetchTasks();
  
  addTaskBtn.addEventListener('click', () => {
    const taskDescription = taskInput.value;
    if (taskDescription) {
      addTask(taskDescription);
      taskInput.value = ''; 
    }
  });

  function fetchTasks() {
    fetch('/tasks', { method: 'GET' })
    .then(response => response.json())
    .then(tasks => {
        taskList.innerHTML = ''; 
        tasks.forEach(addTaskToDOM); 
      })
      .catch(error => console.error('Error fetching tasks: ', error));
  }

  function addTask(description) {
    fetch('/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ description }),
    })
    .then(response => response.json())
    .then(addTaskToDOM) 
    .catch(error => console.error('Error adding task: ', error));
  }

  function addTaskToDOM(task) {
    const li = document.createElement('li');
    li.textContent = task.description;

    // Checkbox for marking task as complete
    const checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.checked = task.completed;
    checkbox.addEventListener('change', () => {
        markTaskComplete(task.id, checkbox.checked);
    });

    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Delete';
    deleteBtn.addEventListener('click', () => {
        deleteTask(task.id);
    });

    li.prepend(checkbox);
    li.appendChild(deleteBtn);
    taskList.appendChild(li);
  }

  function markTaskComplete(taskId, completed) {
    fetch(`/tasks/${taskId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ completed }),
    })
    .then(() => {
        console.log(`Task ${taskId} marked as ${completed ? 'complete' : 'incomplete'}.`);
        fetchTasks();
    })
    .catch(error => console.error('Error updating task:', error));
  }

  function deleteTask(taskId) {
    fetch(`/tasks/${taskId}`, { method: 'DELETE' })
    .then(() => {
        console.log(`Task ${taskId} deleted.`);
        fetchTasks(); 
    })
    .catch(error => console.error('Error deleting task:', error));
  }
});

