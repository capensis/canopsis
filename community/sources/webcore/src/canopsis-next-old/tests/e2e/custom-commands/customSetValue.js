// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customSetValue(selector, value) {
  this.waitForElementVisible(selector)
    .click(selector)
    .setValue(selector, value)
    .pause(WAIT_PAUSE);

  return this;
};
