// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customClearRTE(selector) {
  const { RIGHT_ARROW, BACK_SPACE } = this.Keys;

  this.waitForElementVisible(selector)
    .click(selector)
    .getText(selector, (result) => {
      const chars = result.value.split('');

      chars.forEach(() => this.sendKeys(selector, RIGHT_ARROW));
      chars.forEach(() => this.sendKeys(selector, BACK_SPACE));
    })
    .pause(WAIT_PAUSE);

  return this;
};
