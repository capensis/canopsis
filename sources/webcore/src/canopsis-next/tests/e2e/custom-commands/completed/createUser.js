// http://nightwatchjs.org/guide#usage

module.exports.command = function createUser(
  user = {
    username: 'test-name',
    firstname: 'test-firstname',
    lastname: 'test-lastname',
    email: 'test-email@example.com',
    password: 'test-password',
    role: 1,
    language: 2,
    navigationType: 1,
  },
  callback = result => result,
) {
  const {
    username = 'test-name',
    firstname = 'test-firstname',
    lastname = 'test-lastname',
    email = 'test-email@example.com',
    password = 'test-password',
    role = 1,
    language = 2,
    navigationType = 1,
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
    '/rest/default_rights/user',
    1000,
    () => {
      adminCreateUser.clickSubmitButton();
    },
    (xhr) => {
      const userResponseData = JSON.parse(xhr.responseData)
        .data.filter(item => item.id === username)[0];
      callback({
        userResponseData,
        user,
      });
    },
  );
  adminCreateUser.verifyModalClosed();

  return this;
};
