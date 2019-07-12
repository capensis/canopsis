// http://nightwatchjs.org/guide#usage

module.exports.command = function login(username, password, checkLogin) {
  this.page.auth.login()
    .navigate()
    .verifyPageElementsBefore()
    .clearUsername()
    .setUsername(username)
    .clearPassword()
    .setPassword(password)
    .clickSubmitButton()
    .verifyPageElementsAfter(checkLogin);

  return this;
};
