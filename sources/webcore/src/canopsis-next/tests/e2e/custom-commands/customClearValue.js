// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customClearValue(selector) {
  const { RIGHT_ARROW, BACK_SPACE } = this.Keys;

  this.waitForElementVisible(selector)
    .click(selector)
    .getValue(selector, (result) => {
      const chars = result.value.split('');

      chars.forEach(() => this.setValue(selector, RIGHT_ARROW));
      chars.forEach(() => this.setValue(selector, BACK_SPACE));
    })
    .pause(WAIT_PAUSE);

  return this;
};
