const fs = require('fs');
const path = require('path');

// eslint-disable-next-line import/no-extraneous-dependencies
const { ChartJSNodeCanvas } = require('chartjs-node-canvas');

const chartsFolderPath = path.resolve(process.cwd(), 'benchmarks', '__charts__');

const groupData = (arrayData, prepare) => arrayData.reduce((acc, { name, data }) => {
  acc[name] = prepare ? prepare(data) : data;

  return acc;
}, {});

class ChartsReporter {
  constructor({ width, height }) {
    this.service = new ChartJSNodeCanvas({ width, height, chartCallback: this.chartCallback });
  }

  // eslint-disable-next-line class-methods-use-this
  chartCallback(ChartJS) {
    // eslint-disable-next-line no-param-reassign
    ChartJS.defaults.responsive = true;
    // eslint-disable-next-line no-param-reassign
    ChartJS.defaults.maintainAspectRatio = false;
  }

  generateChart(name, configuration) {
    const buffer = this.service.renderToBufferSync(configuration);

    fs.writeFileSync(path.resolve(chartsFolderPath, `${name}.png`), buffer, 'base64');
  }

  report(...metrics) {
    const groupedMetrics = groupData(
      metrics,
      benchmarkData => groupData(
        benchmarkData,
        groupData,
      ),
    );

    const {
      allProperties,
      allMeasures,
      allBenchmarks,
      allFileNames,
    } = Object.entries(groupedMetrics).reduce((acc, [fileName, benchmarksByName]) => {
      acc.allFileNames.push(fileName);

      Object.entries(benchmarksByName).forEach(([benchmarkName, measuresByName]) => {
        if (!acc.allBenchmarks.includes(benchmarkName)) {
          acc.allBenchmarks.push(benchmarkName);
        }

        Object.entries(measuresByName).forEach(([measureName, properties]) => {
          if (!acc.allMeasures.includes(measureName)) {
            acc.allMeasures.push(measureName);
          }

          if (!properties) {
            return;
          }

          Object.keys(properties).forEach((property) => {
            if (!acc.allProperties.includes(property)) {
              acc.allProperties.push(property);
            }
          });
        });
      });

      return acc;
    }, {
      allProperties: [],
      allMeasures: [],
      allBenchmarks: [],
      allFileNames: [],
    });

    const chartsOptions = allBenchmarks.reduce((acc, benchmarkName) => {
      allProperties.forEach((property) => {
        acc.push({
          name: `${benchmarkName} (${property})`,
          labels: allMeasures,
          datasets: allFileNames.map(fileName => ({
            label: fileName,
            backgroundColor: 'rgba(47,171,99,0.33)',
            data: allMeasures.map(
              measureName => groupedMetrics?.[fileName]?.[benchmarkName]?.[measureName]?.[property] ?? 0,
            ),
          }), []),
        });
      });

      return acc;
    }, []);

    chartsOptions.forEach(({ name, labels, datasets }) => {
      const configuration = {
        type: 'bar',
        data: {
          labels,
          datasets,
        },
        options: {
          plugins: {
            title: {
              display: true,
              text: name,
            },
          },
        },
      };

      this.generateChart(name, configuration);
    });
  }
}

module.exports = ChartsReporter;
