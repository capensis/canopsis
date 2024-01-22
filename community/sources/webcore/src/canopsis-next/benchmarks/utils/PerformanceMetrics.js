/**
 * https://docs.google.com/document/d/1CvAClvFfyA5R-PhYUmn5OOQtYMH4h6I0nSsKchNAySU/preview#heading=h.uxpopqvbjezh
 */
class PerformanceMetrics {
  constructor(data) {
    this.data = JSON.parse(data.toString());
  }

  findLongestPerformanceTask() {
    const allAnimationTasks = this.data.traceEvents.filter(({ name, ph }) => name === 'LongAnimationFrame'
      && ph === 'b');

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
