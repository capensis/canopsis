// A custom Nightwatch assertion.
// The assertion name is the filename.
// Example usage:
//
//   browser.assert.elementsCount(selector, count)
//
// For more information on custom assertions see:
// http://nightwatchjs.org/guide#writing-custom-assertions

module.exports.assertion = function elementsCount(selector, count) {
  this.message = `Testing if element <${selector}> has count: ${count}`;
  this.expected = count;
  this.pass = val => val === count;
  this.value = res => res.value;
  function evaluator(_selector) {
    return document.querySelectorAll(_selector).length;
  }
  this.command = cb => this.api.execute(evaluator, [selector], cb);
};
