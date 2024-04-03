const fs = require('fs');
const path = require('path');

// eslint-disable-next-line import/no-extraneous-dependencies
const { ChartJSNodeCanvas } = require('chartjs-node-canvas');

const { getMetricsGroupedData } = require('../metrics-files');
const { logInfo } = require('../logger');

const CHART_MEASURES_COLORS = [
  '#fda701',
  '#fd693b',
  '#7bb242',
  '#d64315',
  '#fdef75',
  '#fd5252',
  '#9b27af',
];

const getColorByIndex = index => CHART_MEASURES_COLORS[index % CHART_MEASURES_COLORS.length];

class ChartsReporter {
  constructor({ width, height, outputPath }) {
    this.outputPath = outputPath;
    this.service = new ChartJSNodeCanvas({
      width,
      height,
      chartCallback: this.chartCallback,
      backgroundColour: '#fff',
    });
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

    if (!fs.existsSync(this.outputPath)) {
      fs.mkdirSync(this.outputPath);
      logInfo(`Create folder: ${this.outputPath}`);
    }

    const filePath = path.resolve(this.outputPath, `${name}.png`);

    fs.writeFileSync(filePath, buffer, 'base64');
    logInfo(`Save file: ${filePath}`);
  }

  report(...metrics) {
    const { groupedMetrics, allFileNames, allProperties, allMeasures, allBenchmarks } = getMetricsGroupedData(
      metrics,
    );

    const chartsOptions = allBenchmarks.reduce((acc, benchmarkName) => {
      allProperties.forEach((property) => {
        acc.push({
          name: `${benchmarkName} (${property})`,
          labels: allMeasures,
          datasets: allFileNames.map((fileName, index) => ({
            label: fileName,
            backgroundColor: getColorByIndex(index),
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
