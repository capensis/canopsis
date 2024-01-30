const { performance } = require('perf_hooks');

const { enhanceBenchmarkFunction } = require('./iterate');
const Measure = require('./Measure');
const { logInfo } = require('./logger');

class Benchmark {
  constructor(name, benchmarkFunc) {
    this.name = name;
    this.benchmarkFunc = benchmarkFunc;

    this.measuresQueue = [];

    enhanceBenchmarkFunction(this.measure);
  }

  measure = (measureName, measureFunc) => {
    this.measuresQueue.push(new Measure(measureName, measureFunc));
  };

  async run(options) {
    const startTime = performance.now();
    logInfo(`Start benchmark: ${this.name}`);

    await this.benchmarkFunc(this.measure, options);

    await this.measuresQueue.reduce((acc, measure) => acc.then(() => measure.run(options)), Promise.resolve());

    const time = performance.now() - startTime;
    logInfo(`Finish benchmark: ${this.name} (${time.toFixed()} ms)\n`);
  }

  toJSON() {
    return {
      name: this.name,
      data: this.measuresQueue.map(measure => measure.toJSON()),
    };
  }
}

module.exports = Benchmark;
