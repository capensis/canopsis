// http://nightwatchjs.org/guide#usage

module.exports.command = function loginDisabledUser(username, password) {
  this.page.auth.login()
    .navigate()
    .verifyPageElementsBefore()
    .clearUsername()
    .setUsername(username)
    .clearPassword()
    .setPassword(password)
    .clickSubmitButton()
    .verifyErrorDisabledUser();

  return this;
};
