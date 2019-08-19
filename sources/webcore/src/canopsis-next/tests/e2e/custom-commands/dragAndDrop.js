// http://nightwatchjs.org/guide#usage

const { DEFAULT_PAUSE } = require('../config');

module.exports.command = function dragAndDrop(from, to) {
  this
    .moveToElement(from, 5, 5)
    .mouseButtonDown(0)
    .moveToElement(to, 5, 5)
    .mouseButtonUp(0)
    .pause(DEFAULT_PAUSE);

  return this;
};
