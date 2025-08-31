import { defineStore } from 'pinia'

export const useTaskStore = defineStore('task', {
  state: () => ({
    tasks: [],
    currentTask: null
  }),
  actions: {
    setTasks(tasks) {
      this.tasks = tasks
    },
    setCurrentTask(task) {
      this.currentTask = task
    }
  }
})
