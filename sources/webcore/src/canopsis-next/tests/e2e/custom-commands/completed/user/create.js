// http://nightwatchjs.org/guide#usage
const qs = require('qs');
const { USER } = require('../../../constants');
const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createUser(user = { ...USER }, callback = () => {}) {
  const {
    username,
    firstname,
    lastname,
    email,
    password,
    role,
    language,
    navigationType,
  } = user;

  const adminUsersPage = this.page.admin.users();
  const adminCreateUser = this.page.modals.admin.createUser();

  this.url(({ value }) => {
    if (value !== adminUsersPage.url()) {
      adminUsersPage.navigate();
    }
  });

  adminUsersPage.verifyPageElementsBefore()
    .clickAddButton();

  adminCreateUser.verifyModalOpened()
    .setUsername(username)
    .setFirstName(firstname)
    .setLastName(lastname)
    .setEmail(email)
    .setPassword(password)
    .selectRole(role)
    .selectLanguage(language)
    .selectNavigationType(navigationType);

  this.waitForFirstXHR(
    API_ROUTES.user.create,
    1000,
    () => adminCreateUser.clickSubmitButton(),
    ({ responseData, requestData }) => {
      const requestParsedData = qs.parse(requestData);

      return callback({ ...JSON.parse(requestParsedData.user), ...JSON.parse(responseData) });
    },
  );
  adminCreateUser.verifyModalClosed();

  return this;
};
