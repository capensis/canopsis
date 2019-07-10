// http://nightwatchjs.org/guide#usage

const { USERS, CONFIGS } = require('../../constants');

const TEMPORARY_DATA = {};

const onCreateUser = (browser, {
  username, firstname, lastname, email, password,
}) => {
  browser.page.admin.users()
    .clickAddButton();

  browser.page.modals.admin.createUser()
    .verifyModalOpened()
    .setUsername(username)
    .setFirstName(firstname)
    .setLastName(lastname)
    .setEmail(email)
    .setPassword(password)
    .selectRole()
    .selectLanguage()
    .clickSubmitButton()
    .verifyModalClosed();
};

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

  'Create new user with some name': (browser) => {
    const { text, create: { prefix } } = USERS;

    TEMPORARY_DATA[prefix] = {
      username: `${prefix}-${text}-name`,
      firstname: `${prefix}-${text}-firstname`,
      lastname: `${prefix}-${text}-lastname`,
      email: `${prefix}-${text}-email@example.com`,
      password: `${prefix}-${text}-password`,
    };

    browser.page.admin.users()
      .navigate()
      .verifyPageElementsBefore();

    onCreateUser(browser, TEMPORARY_DATA[prefix]);
  },

  'Login new user with some name': async (browser) => {
    const { create: { prefix } } = USERS;
    const { username, password } = TEMPORARY_DATA[prefix];

    await browser.completed.logout()
      .maximizeWindow()
      .completed.login(username, password);
  },

  'Login root': async (browser) => {
    await browser.completed.logout()
      .maximizeWindow()
      .completed.loginAsAdmin();
  },

  'Edit user with some username': (browser) => {
    const { text, create, edit: { prefix } } = USERS;

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
      .clickEditButton(user);

    browser.page.modals.admin.editUser()
      .verifyModalOpened()
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
      .verifyModalClosed();
  },

  'Remove user with some username': (browser) => {
    const { create, edit } = USERS;
    const createUser = TEMPORARY_DATA[create.prefix].username;
    const editUser = TEMPORARY_DATA[edit.prefix].username;

    browser.page.admin.users()
      .verifyPageUserBefore(createUser)
      .clickDeleteButton(createUser)
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .api.pause(CONFIGS.pause);

    browser.page.admin.users()
      .verifyPageUserBefore(editUser)
      .clickDeleteButton(editUser)
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .api.pause(CONFIGS.pause);
  },

  'Create mass users with some name': (browser) => {
    const { text, counts, mass: { prefix } } = USERS;

    TEMPORARY_DATA[prefix] = [];

    for (let i = 0; i < counts; i += 1) {
      TEMPORARY_DATA[prefix].push({
        username: `${prefix}-${text}-name-${i}`,
        firstname: `${prefix}-${text}-firstname-${i}`,
        lastname: `${prefix}-${text}-lastname-${i}`,
        email: `${prefix}-${text}-${i}-email@example.com`,
        password: `${prefix}-${text}-password-${i}`,
      });
    }

    TEMPORARY_DATA[prefix].map(user => onCreateUser(browser, user));
  },

  'Check pagination users table': (browser) => {
    browser.page.admin.users()
      .clickPrevButton()
      .api.pause(CONFIGS.pause);

    browser.page.admin.users()
      .clickNextButton()
      .api.pause(CONFIGS.pause);
  },

  'Delete mass users with some name': (browser) => {
    const { mass: { prefix } } = USERS;

    browser.page.admin.users()
      .selectRange()
      .api.pause(CONFIGS.pause);

    TEMPORARY_DATA[prefix].map(user => browser.page.admin.users()
      .verifyPageUserBefore(user.username)
      .clickOptionCheckbox(user.username)
      .api.pause(CONFIGS.pause));

    browser.page.admin.users()
      .verifyMassDeleteButton()
      .clickMassDeleteButton()
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .api.pause(CONFIGS.pause);
  },
};
