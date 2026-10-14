// http://nightwatchjs.org/guide#usage

const { DEFAULT_PAUSE } = require('../config');

// Default offset for move outside element
const DEFAULT_OFFSET = -5;

/**
 * @param selector - css selector for element
 * @param xOffset - x offset of the element
 * @param yOffset - y offset of the element
 */
module.exports.command = function customClickOutside(
  selector,
  xOffset = DEFAULT_OFFSET,
  yOffset = DEFAULT_OFFSET,
) {
  this.moveTo(selector, xOffset, yOffset)
    .mouseButtonClick(0)
    .pause(DEFAULT_PAUSE);

  return this;
};
