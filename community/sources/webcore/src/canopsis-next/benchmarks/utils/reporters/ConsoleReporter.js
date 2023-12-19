class ConsoleReporter {
  // eslint-disable-next-line class-methods-use-this
  report(benchmarksData) {
    benchmarksData.forEach(({ name, data: benchmarkData }) => {
      const tableData = benchmarkData.map(({ name: measureName, data: measureData }) => ({
        name: measureName,
        ...measureData,
      }));

      /* eslint-disable no-console */
      console.info(name);
      console.table(tableData);
    });
  }
}

module.exports = ConsoleReporter;
