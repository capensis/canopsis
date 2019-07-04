// http://nightwatchjs.org/guide#usage
const { login, users } = require('../../constants');


module.exports = {
  async before(browser, done) {
    const { username, password } = login;

    await browser.maximizeWindow()
      .completed.login(username, password);

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Create new user with some name': (browser) => {
    const {
      username, firstname, lastname, email, password,
    } = users.create;

    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore()
      .clickAddButton()
      .verifyCreateUserModal()
      .setUsername(username)
      .setFirstName(firstname)
      .setLastName(lastname)
      .setEmail(email)
      .setPassword(password)
      .selectRole()
      .clickSubmitButton();
  },

  'Edit user with some username': (browser) => {
    const {
      username, firstname, lastname, email, password,
    } = users.edit;

    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore()
      .clickEditButton()
      .verifyCreateUserModal()
      .clearUsername()
      .setUsername(username)
      .setFirstName(firstname)
      .setLastName(lastname)
      .setEmail(email)
      .setPassword(password)
      .selectRole()
      .clickSubmitButton();
  },

  'Remove user with some username': (browser) => {
    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore()
      .clickDeleteButton()
      .verifyCreateConfirmModal()
      .clickConfirmButton();
  },
};
