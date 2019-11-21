// http://nightwatchjs.org/guide#usage

const { DEFAULT_PAUSE } = require('../config');

module.exports.command = function customClickOutside(selector) {
  this.moveTo(selector, -5, -5)
    .mouseButtonClick(0)
    .pause(DEFAULT_PAUSE);

  return this;
};
