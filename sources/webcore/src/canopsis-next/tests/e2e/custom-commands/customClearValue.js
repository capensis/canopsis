// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customClearValue(selector) {
  const { CONTROL, DELETE } = this.Keys;

  this.customClick(selector)
    .sendKeys(selector, [CONTROL, 'a'])
    .sendKeys(selector, DELETE)
    .pause(WAIT_PAUSE);

  return this;
};
