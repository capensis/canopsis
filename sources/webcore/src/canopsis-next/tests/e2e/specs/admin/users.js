// http://nightwatchjs.org/guide#usage

const v = 'asd';

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .login('root', 'root');

    done();
  },

  after(browser, done) {
    browser.end(done);
  },
  'Create new user with some name': (browser) => {
    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore()
      .verifyCreateUserModal()
      .clickAddButton()
      .setUsername(v)
      .setFirstName(v)
      .setLastName(v)
      .setEmail('asd@asd.asd')
      .setPassword(v)
      .selectRole();
  },

  // 'Edit user with some username': (browser) => {},
  //
  // 'Remove user with some username': (browser) => {},
};
