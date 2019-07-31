// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../src/config');

module.exports.command = function deleteUser(id, callback = () => {}) {
  const adminUsersPage = this.page.admin.users();
  const confirmation = this.page.modals.confirmation();

  adminUsersPage.navigate()
    .selectRange()
    .verifyPageUserBefore(id)
    .clickDeleteButton(id);

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    `${API_ROUTES.user.remove}/${id}`,
    5000,
    () => confirmation.clickConfirmButton(),
    ({ responseData }) => callback(JSON.parse(responseData)),
  );

  confirmation.verifyModalClosed();

  return this;
};
