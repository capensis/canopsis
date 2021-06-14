// http://nightwatchjs.org/guide#usage

module.exports.command = function customClick(selector) {
  this.useXpath()
    .customClick(selector)
    .useCss();

  return this;
};
