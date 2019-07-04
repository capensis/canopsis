// http://nightwatchjs.org/guide#usage
const { login, users } = require('../../constants');

const WAIT_PAUSE = 500;

module.exports = {
  async before(browser, done) {
    const { username, password } = login;

    await browser.maximizeWindow()
      .completed.login(username, password);

    done();
  },

  async after(browser, done) {
    await browser.completed.logout()
      .end(done);
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
      .clickSubmitButton()
      .api.pause(WAIT_PAUSE);
  },

  'Edit user with some username': (browser) => {
    const {
      username, firstname, lastname, email, password,
    } = users.edit;

    browser.page.admin.users()
      .clickEditButton()
      .verifyCreateUserModal()
      .clearUsername()
      .setUsername(username)
      .clearFirstName()
      .setFirstName(firstname)
      .clearLastName()
      .setLastName(lastname)
      .clearEmail()
      .setEmail(email)
      .setPassword(password)
      .selectRole(4)
      .clickSubmitButton()
      .api.pause(WAIT_PAUSE);
  },

  'Remove user with some username': (browser) => {
    browser.page.admin.users()
      .clickDeleteButton()
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .api.pause(WAIT_PAUSE);
  },
};
