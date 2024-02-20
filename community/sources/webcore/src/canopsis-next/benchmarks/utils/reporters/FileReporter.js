const fs = require('fs');
const path = require('path');

class FileReporter {
  static readMetricFile(outputPath, metricName) {
    const fileName = metricName.endsWith('.json') ? metricName : `${metricName}.json`;

    if (!fs.existsSync(path.resolve(outputPath, fileName))) {
      throw Error(`Metric "${fileName}" doesn't exist`);
    }

    return {
      name: metricName.replace('.json', ''),
      data: JSON.parse(fs.readFileSync(path.resolve(outputPath, fileName)).toString()),
    };
  }

  static readMetricFiles(outputPath) {
    if (fs.existsSync(outputPath)) {
      const files = fs.readdirSync(outputPath);

      return files
        .filter(name => name.endsWith('.json'))
        .map(name => FileReporter.readMetricFile(outputPath, name));
    }

    return [];
  }

  constructor({ name = 'metrics', outputPath } = {}) {
    this.outputPath = outputPath;
    this.name = name;
  }

  // eslint-disable-next-line class-methods-use-this
  report(data) {
    if (!fs.existsSync(this.outputPath)) {
      fs.mkdirSync(this.outputPath);
    }

    const jsonName = this.name.endsWith('.json') ? this.name : `${this.name}.json`;

    fs.writeFileSync(path.resolve(this.outputPath, jsonName), JSON.stringify(data, undefined, 2));
  }
}

module.exports = FileReporter;
