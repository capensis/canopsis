// http://nightwatchjs.org/guide#usage

const { USERS } = require('../../constants');

const TEMPORARY_DATA = {};

const onCreateTemporaryObject = ({ prefix, text, index }) => {
  const i = typeof index === 'number' ? `-${index}` : '';
  return {
    username: `${prefix}-${text}${i}-name`,
    firstname: `${prefix}-${text}${i}-firstname`,
    lastname: `${prefix}-${text}${i}-lastname`,
    email: `${prefix}-${text}${i}-email@example.com`,
    password: `${prefix}-${text}${i}-password`,
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

    browser.completed.createUser(TEMPORARY_DATA[prefix], ({ userResponseData }) => {
      TEMPORARY_DATA[prefix] = {
        ...TEMPORARY_DATA[prefix],
        userResponseData,
      };
    });
  },

  'Login by created user credentials': (browser) => {
    const { create: { prefix } } = USERS;
    const { username, password } = TEMPORARY_DATA[prefix];

    browser.completed.logout()
      .maximizeWindow()
      .completed.login(username, password);
  },

  'Edit special user with username from constants': (browser) => {
    browser.completed.logout()
      .maximizeWindow()
      .completed.loginAsAdmin();

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

    browser.page.modals.admin.createUser()
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
      .clickEnabled()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Login by disabled user credentials': (browser) => {
    const { edit: { prefix } } = USERS;
    const { username, password } = TEMPORARY_DATA[prefix];

    browser.completed.logout()
      .maximizeWindow()
      .completed.loginDisabledUser(username, password);
  },

  'Remove user with some username': (browser) => {
    browser.completed.loginAsAdmin();

    const { create, edit } = USERS;
    const createUser = TEMPORARY_DATA[create.prefix].username;
    const editUser = TEMPORARY_DATA[edit.prefix].username;

    browser.completed.deleteUser({
      username: createUser,
    }, ({ responseData }) => {
      TEMPORARY_DATA[create.prefix] = {
        ...TEMPORARY_DATA[create.prefix],
        deleteData: responseData,
      };
    });

    browser.completed.deleteUser({
      username: editUser,
    }, ({ responseData }) => {
      TEMPORARY_DATA[edit.prefix] = {
        ...TEMPORARY_DATA[edit.prefix],
        deleteData: responseData,
      };
    });
  },

  'Create mass users with some name': (browser) => {
    const { text, counts, mass: { prefix } } = USERS;

    TEMPORARY_DATA[prefix] = [];

    for (let index = 0; index < counts; index += 1) {
      TEMPORARY_DATA[prefix].push(onCreateTemporaryObject({ text, prefix, index }));
    }

    TEMPORARY_DATA[prefix].map((user, index) => browser.completed.createUser(user, ({ userResponseData }) => {
      TEMPORARY_DATA[prefix][index] = {
        ...TEMPORARY_DATA[prefix][index],
        userResponseData,
      };
    }));
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
      .clickMassDeleteButton();
    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();
  },
};
