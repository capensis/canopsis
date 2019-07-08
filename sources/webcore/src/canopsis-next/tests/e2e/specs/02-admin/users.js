// http://nightwatchjs.org/guide#usage
const { login, users } = require('../../constants');

const WAIT_PAUSE = 500;

const TEMPORARY_DATA = {};

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
    const { text, create: { prefix } } = users;

    TEMPORARY_DATA[prefix] = {
      username: `${prefix}-${text}-name`,
      firstname: `${prefix}-${text}-firstname`,
      lastname: `${prefix}-${text}-lastname`,
      email: `${prefix}-${text}-email@example.com`,
      password: `${prefix}-${text}-password`,
    };

    const {
      username, firstname, lastname, email, password,
    } = TEMPORARY_DATA[prefix];

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
      .selectLanguage()
      .clickSubmitButton()
      .api.pause(WAIT_PAUSE);
  },

  'Login new user with some name': async (browser) => {
    const { create: { prefix } } = users;
    const { username, password } = TEMPORARY_DATA[prefix];

    await browser.completed.logout()
      .maximizeWindow()
      .completed.login(username, password);
  },

  'Login root': async (browser) => {
    const { username, password } = login;

    await browser.completed.logout()
      .maximizeWindow()
      .completed.login(username, password);
  },

  'Edit user with some username': (browser) => {
    const { text, create, edit: { prefix } } = users;

    TEMPORARY_DATA[prefix] = {
      username: `${prefix}-${text}-name`,
      firstname: `${prefix}-${text}-firstname`,
      lastname: `${prefix}-${text}-lastname`,
      email: `${prefix}-${text}-email@example.com`,
      password: `${prefix}-${text}-password`,
    };

    const user = TEMPORARY_DATA[create.prefix].username;

    const {
      username, firstname, lastname, email, password,
    } = TEMPORARY_DATA[prefix];

    browser.page.admin.users()
      .navigate()
      .verifyPageUserBefore(user)
      .clickEditButton(user)
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
      .selectLanguage(2)
      .clickSubmitButton()
      .api.pause(WAIT_PAUSE);
  },

  'Remove user with some username': (browser) => {
    const { create, edit } = users;
    const createUser = TEMPORARY_DATA[create.prefix].username;
    const editUser = TEMPORARY_DATA[edit.prefix].username;

    browser.page.admin.users()
      .verifyPageUserBefore(editUser)
      .clickDeleteButton(editUser)
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .api.pause(WAIT_PAUSE)
      .verifyPageUserBefore(createUser)
      .clickDeleteButton(createUser)
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .api.pause(WAIT_PAUSE);
  },
};
