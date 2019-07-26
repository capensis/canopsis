// http://nightwatchjs.org/guide#usage
const { USER } = require('../../constants');
const { API_ROUTES: { user: { remove } } } = require('../../../../src/config');

module.exports.command = function deleteUser(
  user = USER,
  callback = result => result,
) {
  const {
    username = USER.username,
  } = user;

  const adminUsersPage = this.page.admin.users();
  const confirmation = this.page.modals.confirmation();

  adminUsersPage.navigate()
    .selectRange()
    .verifyPageUserBefore(username)
    .clickDeleteButton(username);

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    `${remove}/${username}`,
    1000,
    () => confirmation.clickConfirmButton(),
    ({ responseData }) => {
      callback(JSON.parse(responseData));
    },
  );

  confirmation.verifyModalClosed();

  return this;
};
