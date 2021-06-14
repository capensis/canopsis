// http://nightwatchjs.org/guide#usage

module.exports.command = function logout() {
  this.page.auth.logout()
    .verifyPageElementsBefore();

  this.page.layout.popup()
    .clickOnEveryPopupsCloseIcons();

  this.page.layout.topBar()
    .clickUserDropdown()
    .clickLogoutButton();

  this.page.auth.logout()
    .verifyPageElementsAfter();

  return this;
};
