// http://nightwatchjs.org/guide#usage

const WAIT_PAUSE = 500;

module.exports.command = function customClearRTE(selector) {
  const { BACK_SPACE, END } = this.Keys;

  this.customClick(selector)
    .sendKeys(selector, [END])
    .getText(selector, (result) => {
      this.sendKeys(selector, new Array(result.value.length).fill(BACK_SPACE));
    })
    .pause(WAIT_PAUSE);

  return this;
};
