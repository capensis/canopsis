// http://nightwatchjs.org/guide#usage

const { DEFAULT_PAUSE } = require('../config');

module.exports.command = function customClick(selector) {
  this.waitForElementVisible(selector)
    .click(selector)
    .pause(DEFAULT_PAUSE);

  return this;
};
