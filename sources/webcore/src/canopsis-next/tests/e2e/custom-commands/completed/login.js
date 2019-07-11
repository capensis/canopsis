// http://nightwatchjs.org/guide#usage

module.exports.command = function login(username, password) {
  this.page.auth.login()
    .navigate()
    .verifyPageElementsBefore()
    .setUsername(username)
    .setPassword(password)
    .clickSubmitButton()
    .verifyPageElementsAfter();

  return this;
};
