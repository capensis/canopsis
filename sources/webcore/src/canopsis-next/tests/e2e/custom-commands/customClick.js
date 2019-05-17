// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customClick(selector) {
  this.waitForElementVisible(selector)
    .click(selector)
    .pause(WAIT_PAUSE);

  return this;
};
