// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customKeyup(selector, key) {
  this.waitForElementVisible(selector)
    .setValue(selector, this.Keys[key])
    .pause(WAIT_PAUSE);

  return this;
};
