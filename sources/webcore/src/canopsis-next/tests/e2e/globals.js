// http://nightwatchjs.org/guide#external-globals

const nightWatchRecord = require('nightwatch-record');

module.exports = {
  asyncHookTimeout: 20000,
  waitForConditionTimeout: 5000,

  beforeEach(browser, done) {
    nightWatchRecord.start(browser, done);
  },

  afterEach(browser, done) {
    nightWatchRecord.stop(browser, done);
  },
};
