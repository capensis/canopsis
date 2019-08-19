// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  async after(browser, done) {
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
      .clearFooter()
      .setFooter('TestFooter')
      .clearDescription()
      .setDescription('TestDescription')
      .clickSubmitButton();
  },
};
