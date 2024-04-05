/**
 * https://docs.google.com/document/d/1CvAClvFfyA5R-PhYUmn5OOQtYMH4h6I0nSsKchNAySU/preview#heading=h.uxpopqvbjezh
 */
class PerformanceMetrics {
  constructor(data) {
    this.data = JSON.parse(data.toString());
  }

  filterTasks(callback) {
    return this.data.traceEvents.filter(callback);
  }

  getTasksByName(name) {
    return this.filterTasks(task => task.name === name);
  }

  findFinishXHRTaskByUrl(url) {
    return this.getTasksByName('XHRReadyStateChange').find(
      task => task.args.data.url.includes(url) && task.args.data.readyState === 4,
    );
  }

  findLongestPerformanceTask() {
    const allAnimationTasks = this.getTasksByName('LongAnimationFrame').filter(({ ph }) => ph === 'b');

    return allAnimationTasks.reduce((acc, task) => {
      if (acc.args.data.duration < task.args.data.duration) {
        return task;
      }

      return acc;
    });
  }
}

module.exports = {
  PerformanceMetrics,
};
