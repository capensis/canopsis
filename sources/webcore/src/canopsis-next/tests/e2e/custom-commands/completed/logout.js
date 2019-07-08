// http://nightwatchjs.org/guide#usage

module.exports.command = function logout() {
  this.page.auth.logout()
    .verifyPageElementsBefore();

  this.page.layout()
    .clickOnEveryPopupsCloseIcons();

  this.page.auth.logout()
    .clickUserNavigationTopBarButton()
    .clickLogoutButton();

  return this;
};
