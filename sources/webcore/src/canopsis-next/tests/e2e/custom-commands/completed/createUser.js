// http://nightwatchjs.org/guide#usage
const { USER } = require('../../constants');
const { API_ROUTES: { user: { list } } } = require('../../../../src/config');
const fs = require('fs');

module.exports.command = function createUser(
  user = USER,
  callback = result => fs.writeFile('user.json', JSON.stringify(result), 'utf8', () => {}),
) {
  const {
    username = USER.username,
    firstname = USER.firstname,
    lastname = USER.lastname,
    email = USER.email,
    password = USER.password,
    role = USER.role,
    language = USER.language,
    navigationType = USER.navigationType,
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
    list,
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
