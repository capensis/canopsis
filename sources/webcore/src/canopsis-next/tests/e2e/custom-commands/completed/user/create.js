// http://nightwatchjs.org/guide#usage
const qs = require('qs');
const { USER, WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');
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
    .selectNavigationType(navigationType);

  if (language) {
    adminCreateUser.selectLanguage(language);
  }

  this.waitForFirstXHR(
    API_ROUTES.user.create,
    WAIT_FOR_FIRST_XHR_TIME,
    () => adminCreateUser.clickSubmitButton(),
    ({ responseData, requestData }) => {
      adminCreateUser.verifyModalClosed();

      const requestParsedData = qs.parse(requestData);

      callback({ ...JSON.parse(requestParsedData.user), ...JSON.parse(responseData) });
    },
  );

  return this;
};
