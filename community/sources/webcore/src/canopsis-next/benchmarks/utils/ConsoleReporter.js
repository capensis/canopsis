class ConsoleReporter {
  // eslint-disable-next-line class-methods-use-this
  report(data) {
    Object.entries(data).forEach(([benchmarkName, benchmarkData]) => {
      /* eslint-disable no-console */
      console.info(benchmarkName);
      console.table(benchmarkData);
    });
  }
}

module.exports = ConsoleReporter;
