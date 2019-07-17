// http://nightwatchjs.org/guide#usage

const { DEFAULT_PAUSE } = require('../config');

module.exports.command = function defaultPause() {
  this.pause(DEFAULT_PAUSE);
  return this;
};
