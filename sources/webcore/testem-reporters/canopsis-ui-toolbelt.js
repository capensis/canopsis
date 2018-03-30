var resultString = '| Browser | Test name | Status | Log |\n| ------ | ------ | ------ | --- |\n';

function StandardReporter(silent, out) {
  this.out = out || process.stdout;
  this.silent = silent;
  this.stoppedOnError = null;
  this.id = 1;
  this.total = 0;
  this.pass = 0;
  this.skipped = 0;
  this.results = [];
  this.errors = [];
  this.logs = [];
}
StandardReporter.prototype = {
  report: function(prefix, data) {
    this.results.push({
      launcher: prefix,
      result: data
    });
    this.display(prefix, data);
    this.total++;
    if (data.skipped) {
      this.skipped++;
    } else if (data.passed) {
      this.pass++;
    }
  },
  summaryDisplay: function() {
    // resultDict.total = this.total;
    // resultDict.pass = this.pass;
    // resultDict.skip = this.skip;
    // resultDict.fail = this.total - this.pass - this.skipped;

    // if (this.pass + this.skipped === this.total) {
    //   resultDict.message = 'OK';
    //   resultDict.state = 0;
    // } else {
    //   resultDict.message = 'NOK';
    //   resultDict.state = 2;
    // }
  },
  display: function(prefix, result) {
    if (this.silent) {
      return;
    }
    var dict = result;
    dict.prefix = prefix;
    dict.id = this.id++;
    if(!dict.logs) dict.logs = [];
    dict.logs = dict.logs.filter(function (el) {
      return el.type !== "warn" && el.type !== "log";
    });
    resultString += '| ' + result.prefix + ' | ' + result.name + ' | ' + (result.passed ? ':ok: OK' : ':negative_squared_cross_mark: ERROR') + ' | ' + JSON.stringify(dict.logs) + ' |\n';
  },
  finish: function() {
    if (this.silent) {
      return;
    }
    this.summaryDisplay();
    this.out.write(resultString);
  }
};

module.exports = StandardReporter;
