// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.login('root', 'root'); // TODO: use from some constants file

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Create new user with some name': (browser) => {
    const value = 'asd'; // TODO: use from some constants file

    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore()
      .clickAddButton()
      .verifyCreateUserModal()
      .setUsername(value)
      .setFirstName(value)
      .setLastName(value)
      .setEmail(`${value}@${value}.com`)
      .setPassword(value)
      .selectRole()
      .clickSubmitButton();
  },

  // 'Edit user with some username': (browser) => {},
  //
  // 'Remove user with some username': (browser) => {},
};
