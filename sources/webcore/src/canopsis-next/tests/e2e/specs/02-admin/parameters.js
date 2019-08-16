// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    browser.globals.parameters = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  async after(browser, done) {
    delete browser.globals.parameters;

    await browser.completed.logout()
      .end(done);
  },

  'Change parameters with some name': (browser) => {
    browser.page.admin.parameters()
      .navigate()
      .verifyPageElementsBefore()
      .clearAppTitle()
      .setAppTitle('TestAppTitle')
      .selectLanguage(2)
      .setFooter('TestFooter')
      .setDescription('TestDescription')
      .clickSubmitButton();
  },
};
