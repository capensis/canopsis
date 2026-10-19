// http://nightwatchjs.org/guide#external-globals

const nightWatchRecord = require('nightwatch-record');
const nightWatchRecordConfig = require('./nightwatch-record.config');

module.exports = {
  asyncHookTimeout: 50000,
  waitForConditionTimeout: 5000,
  test_settings: {
    videos: nightWatchRecordConfig,
  },

  beforeEach(browser, done) {
    nightWatchRecord.start(browser, done);
  },

  afterEach(browser, done) {
    nightWatchRecord.stop(browser, done);
  },
};
