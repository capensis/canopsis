// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function deleteUser(id, callback = () => {}) {
  const adminUsersPage = this.page.admin.users();
  const confirmation = this.page.modals.common.confirmation();

  adminUsersPage.navigate()
    .selectRange()
    .verifyPageUserBefore(id)
    .clickDeleteButton(id);

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    `${API_ROUTES.user.remove}/${id}`,
    5000,
    () => confirmation.clickSubmitButton(),
    ({ responseData }) => {
      confirmation.verifyModalClosed();

      callback(JSON.parse(responseData));
    },
  );

  return this;
};
