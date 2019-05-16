// For authoring Nightwatch tests, see
// http://nightwatchjs.org/guide#usage

module.exports.command = function logout() {
  this.page.auth.logout()
    .navigate()
    .clickUserNavigationTopBarButton()
    .clickLogoutButton();

  return this;
};
