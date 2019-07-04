// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customClearValue(selector) {
  this.waitForElementVisible(selector)
    .click(selector)
    .clearValue(selector)
    .expect.element(selector).text.to.equal('')
    .pause(WAIT_PAUSE);

  return this;
};
