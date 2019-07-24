// http://nightwatchjs.org/guide#usage

module.exports.command = function deleteUser(
  user = {
    username: 'test-name',
  },
  callback = result => result,
) {
  const {
    username = 'test-name',
  } = user;

  const adminUsersPage = this.page.admin.users();
  const confirmation = this.page.modals.confirmation();

  adminUsersPage.navigate()
    .selectRange()
    .verifyPageUserBefore(username)
    .clickDeleteButton(username);

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    `/account/delete/user/${username}`,
    1000,
    () => {
      confirmation.clickConfirmButton();
    },
    (xhr) => {
      const responseData = JSON.parse(xhr.responseData);
      callback({ responseData });
    },
  );

  confirmation.verifyModalClosed();

  return this;
};
