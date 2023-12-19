const fs = require('fs');
const path = require('path');

const metricsFolderPath = path.resolve(process.cwd(), 'benchmarks', '__metrics__');

class FileReporter {
  static readMetricFiles() {
    if (fs.existsSync(metricsFolderPath)) {
      const files = fs.readdirSync(metricsFolderPath);

      return files
        .filter(name => name.endsWith('.json'))
        .map(name => ({
          name: name.replace('.json', ''),
          data: JSON.parse(fs.readFileSync(path.resolve(metricsFolderPath, name)).toString()),
        }));
    }

    return [];
  }

  constructor({ name = 'metrics' } = {}) {
    this.name = name;
  }

  // eslint-disable-next-line class-methods-use-this
  report(data) {
    if (!fs.existsSync(metricsFolderPath)) {
      fs.mkdirSync(metricsFolderPath);
    }

    const jsonName = this.name.endsWith('.json') ? this.name : `${this.name}.json`;

    fs.writeFileSync(path.resolve(metricsFolderPath, jsonName), JSON.stringify(data, undefined, 2));
  }
}

module.exports = FileReporter;
