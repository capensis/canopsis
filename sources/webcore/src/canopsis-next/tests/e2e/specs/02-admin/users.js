// http://nightwatchjs.org/guide#usage
const { generateTemporaryUser } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.users = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  async after(browser, done) {
    delete browser.globals.users;

    await browser.completed.logout()
      .end(done);
  },

  'Create new user with some name': (browser) => {
    const { users } = browser.globals;

    users.create = generateTemporaryUser('create');

    browser.completed.createUser(users.create, ({ userResponseData }) => {
      users.create = {
        ...users.create,
        userResponseData,
      };
    });
  },

  'Login by created user credentials': (browser) => {
    const { username, password } = browser.globals.users.create;

    browser.completed.logout()
      .maximizeWindow()
      .completed.login(username, password);
  },

  'Edit special user with username from constants': (browser) => {
    browser.completed.logout()
      .maximizeWindow()
      .completed.loginAsAdmin();

    const { users } = browser.globals;

    users.edit = generateTemporaryUser('edit');

    const userSelector = users.create.username;

    const {
      username, firstname, lastname, email, password,
    } = users.edit;

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
    const { username, password } = browser.globals.users.edit;

    browser.completed.logout()
      .maximizeWindow()
      .completed.loginDisabledUser(username, password);
  },

  'Remove user with some username': (browser) => {
    const { users } = browser.globals;

    browser.completed.loginAsAdmin();

    const createUser = users.create.username;
    const editUser = users.edit.username;

    browser.completed.deleteUser(createUser, ({ responseData }) => {
      users.create = {
        ...users.create,
        deleteData: responseData,
      };
    });

    browser.completed.deleteUser(editUser, ({ responseData }) => {
      users.edit = {
        ...users.edit,
        deleteData: responseData,
      };
    });
  },

  'Create mass users with some name': (browser) => {
    const { users } = browser.globals;

    users.mass = [];

    for (let index = 0; index < 5; index += 1) {
      users.mass.push(generateTemporaryUser('mass'));
    }

    users.mass.map((user, index) => browser.completed.createUser(user, ({ userResponseData }) => {
      users.mass[index] = {
        ...users.mass[index],
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
    const { users } = browser.globals;

    browser.page.admin.users()
      .selectRange()
      .defaultPause();

    users.mass.map(user => browser.page.admin.users()
      .verifyPageUserBefore(user.username)
      .clickOptionCheckbox(user.username)
      .defaultPause());

    browser.page.admin.users()
      .verifyMassDeleteButton()
      .clickMassDeleteButton();
    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },
};
