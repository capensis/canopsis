// http://nightwatchjs.org/guide#usage

const { isNaN } = require('lodash');
const { USERS } = require('../../constants');

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

const onCreateTemporaryObject = ({ prefix, text, i }) => {
  const index = isNaN(i) ? '' : `-${i}`;
  return {
    username: `${prefix}-${text}${index}-name`,
    firstname: `${prefix}-${text}${index}-firstname`,
    lastname: `${prefix}-${text}${index}-lastname`,
    email: `${prefix}-${text}${index}-email@example.com`,
    password: `${prefix}-${text}${index}-password`,
  };
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

    TEMPORARY_DATA[prefix] = onCreateTemporaryObject({ text, prefix });

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

    TEMPORARY_DATA[prefix] = onCreateTemporaryObject({ text, prefix });

    const userSelector = TEMPORARY_DATA[create.prefix].username;

    const {
      username, firstname, lastname, email, password,
    } = TEMPORARY_DATA[prefix];

    browser.page.admin.users()
      .navigate()
      .verifyPageUserBefore(userSelector)
      .clickEditButton(userSelector);

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
      .defaultPause();

    browser.page.admin.users()
      .verifyPageUserBefore(editUser)
      .clickDeleteButton(editUser)
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .defaultPause();
  },

  'Create mass users with some name': (browser) => {
    const { text, counts, mass: { prefix } } = USERS;

    TEMPORARY_DATA[prefix] = [];

    for (let i = 0; i < counts; i += 1) {
      TEMPORARY_DATA[prefix].push(onCreateTemporaryObject({ text, prefix, i }));
    }

    TEMPORARY_DATA[prefix].map(user => onCreateUser(browser, user));
  },

  'Check pagination users table': (browser) => {
    browser.page.admin.users()
      .clickPrevButton()
      .defaultPause();

    browser.page.admin.users()
      .clickNextButton()
      .defaultPause();
  },

  'Delete mass users with some name': (browser) => {
    const { mass: { prefix } } = USERS;

    browser.page.admin.users()
      .selectRange()
      .defaultPause();

    TEMPORARY_DATA[prefix].map(user => browser.page.admin.users()
      .verifyPageUserBefore(user.username)
      .clickOptionCheckbox(user.username)
      .defaultPause());

    browser.page.admin.users()
      .verifyMassDeleteButton()
      .clickMassDeleteButton()
      .verifyCreateConfirmModal()
      .clickConfirmButton()
      .defaultPause();
  },
};
