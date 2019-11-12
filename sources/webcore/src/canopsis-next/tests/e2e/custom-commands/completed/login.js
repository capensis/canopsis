// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../src/config');

module.exports.command = function login(username, password) {
  const loginPage = this.page.auth.login();

  loginPage.navigate()
    .verifyPageElementsBefore()
    .clearUsername()
    .setUsername(username)
    .clearPassword()
    .setPassword(password);

  this.waitForFirstXHR(
    `${API_ROUTES.currentUser}`,
    5000,
    () => loginPage.clickSubmitButton(),
    ({ responseData }) => {
      const { data: [user] } = JSON.parse(responseData);

      this.globals.currentUser = user;

      loginPage.verifyPageElementsAfter();
    },
  );

  return this;
};
