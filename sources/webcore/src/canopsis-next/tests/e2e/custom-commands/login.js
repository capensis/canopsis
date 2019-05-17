// http://nightwatchjs.org/guide#usage

module.exports.command = function login(username, password) {
  this.page.auth.login()
    .navigate()
    .verifyPageElementsBefore()
    .enterUsername(username)
    .enterPassword(password)
    .clickSubmitButton()
    .verifyPageElementsAfter();

  return this;
};
