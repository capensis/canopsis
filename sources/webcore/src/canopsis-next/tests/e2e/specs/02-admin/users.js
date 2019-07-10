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

  'Create new user with some name': (browser) => {
    const value = 'asd'; // TODO: use from some constants file when we will use fixtures

    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore()
      .clickAddButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .setUsername(value)
      .setFirstName(value)
      .setLastName(value)
      .setEmail(`${value}@${value}.com`)
      .setPassword(value)
      .selectRole()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  // 'Edit user with some username': (browser) => {},
  //
  // 'Remove user with some username': (browser) => {},
};
