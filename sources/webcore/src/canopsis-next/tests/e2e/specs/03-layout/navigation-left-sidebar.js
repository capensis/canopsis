// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Open sidebar': (browser) => {
    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .defaultPause();
  },

  'Add view with some name': (browser) => {
    browser.page.layout.leftSideBar()
      .verifySettingsWrapperBefore()
      .clickSettingsViewButton()
      .verifyControlsWrapperBefore()
      .clickAddViewButton()
      .defaultPause();
  },
};
