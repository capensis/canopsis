// http://nightwatchjs.org/guide#usage

module.exports.command = function logout() {
  this.page.auth.logout()
    .verifyPageElementsBefore()
    .clickUserNavigationTopBarButton()
    .clickLogoutButton();

  return this;
};
