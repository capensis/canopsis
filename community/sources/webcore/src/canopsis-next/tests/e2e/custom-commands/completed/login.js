// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../src/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../constants');

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
    WAIT_FOR_FIRST_XHR_TIME,
    () => loginPage.clickSubmitButton(),
    ({ responseData }) => {
      const { data: [user] } = JSON.parse(responseData);

      this.globals.currentUser = user;

      loginPage.verifyPageElementsAfter();
    },
  );

  return this;
};
