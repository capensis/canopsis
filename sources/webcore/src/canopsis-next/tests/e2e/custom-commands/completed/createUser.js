// http://nightwatchjs.org/guide#usage
const { USER } = require('../../constants');
const { API_ROUTES } = require('../../../../src/config');

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

  adminUsersPage.navigate()
    .verifyPageElementsBefore()
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
    API_ROUTES.user.list,
    1000,
    () => adminCreateUser.clickSubmitButton(),
    ({ responseData }) => {
      const userResponseData = JSON.parse(responseData)
        .data.filter(item => item.id === username)[0];
      callback(userResponseData);
    },
  );
  adminCreateUser.verifyModalClosed();

  return this;
};
