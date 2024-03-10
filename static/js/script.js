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
            if (!tasks || !Array.isArray(tasks)) {
                throw new Error('Invalid tasks data');
            }
            taskList.innerHTML = '';
            tasks.forEach(addTaskToDOM);
        })
        .catch(error => console.error('Error fetching tasks:', error));
    }

    function addTask(description) {
    fetch('/tasks', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ Description: description }), // Adjust to match your Go struct tags
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        if (response.status === 204 || response.statusText === 'No Content') {
            // Handle no content response
            return {};
        }
        return response.json(); // This will fail if the response is empty, so ensure it's safe to call
    })
    .then(task => {
        if (task && task.ID) { // Ensure there's a task object and an ID before trying to add it to the DOM
            addTaskToDOM(task);
        }
        fetchTasks(); // Fetch tasks again to refresh the list, even if no task was added to the DOM
    })
    .catch(error => console.error('Error adding task:', error));
}

    function markTaskComplete(taskId, taskDescription, completed) {
        fetch(`/tasks/${taskId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ Completed: completed }), // Adjust to match your Go struct tags
        })
            .then(() => {
        console.log(`Task ${taskId}: ${taskDescription} marked as ${completed ? 'complete' : 'incomplete'}.`);
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

function addTaskToDOM(task) {
    const li = document.createElement('li');
    li.setAttribute('data-id', task.ID);

    const descriptionSpan = document.createElement('span');
    descriptionSpan.textContent = task.Description; 
    li.appendChild(descriptionSpan);

    const checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.checked = task.Completed;
    checkbox.addEventListener('change', function() {
        const taskId = this.parentNode.getAttribute('data-id');
        markTaskComplete(taskId, descriptionSpan.textContent,  this.checked);
    });
    li.appendChild(checkbox);

    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Delete';
    deleteBtn.addEventListener('click', function() {
        const taskId = this.parentNode.getAttribute('data-id');
        deleteTask(taskId);
    });
    li.appendChild(deleteBtn);

    taskList.appendChild(li);
}

});

