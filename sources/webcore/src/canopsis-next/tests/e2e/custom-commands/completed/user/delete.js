// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

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
    WAIT_FOR_FIRST_XHR_TIME,
    () => confirmation.clickSubmitButton(),
    ({ responseData }) => {
      confirmation.verifyModalClosed();

      callback(JSON.parse(responseData));
    },
  );

  return this;
};
