const { performance } = require('perf_hooks');

const { logInfo } = require('./logger');

class Measure {
  constructor(name, measureFunc) {
    this.name = name;
    this.measureFunc = measureFunc;
  }

  report = (data) => {
    this.reportData = data;
  };

  toJSON() {
    return {
      name: this.name,
      data: this.reportData,
    };
  }

  async run(options) {
    const startTime = performance.now();
    logInfo(`Start metric: ${this.name}`);

    await this.measureFunc(this.report, options);

    const time = performance.now() - startTime;

    logInfo(`Finish metric: ${this.name} (${time.toFixed()} ms)`);
  }
}

module.exports = Measure;
