// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../src/config');
const { generateTemporaryUser } = require('../../helpers/entities');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../constants');

module.exports = {
  async before(browser, done) {
    browser.globals.users = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    delete browser.globals.users;

    browser.completed.logout()
      .end(done);
  },

  'Create new user with some name': (browser) => {
    const { users } = browser.globals;
    const generatedUser = generateTemporaryUser();

    browser.completed.user.create(generatedUser, user => users.create = { ...user, ...generatedUser });
  },

  'Check searching': (browser) => {
    const { create: user } = browser.globals.users;
    const usersPage = browser.page.admin.users();

    usersPage.setSearchingText(user._id)
      .waitForFirstXHR(
        API_ROUTES.user.list,
        WAIT_FOR_FIRST_XHR_TIME,
        () => usersPage.clickSubmitSearchButton(), ({ responseData }) => {
          const { data } = JSON.parse(responseData);

          browser.assert.ok(data.every(item => item._id === user._id));
          browser.assert.elementsCount(usersPage.elements.dataTableUserItem.selector, 1);

          usersPage.verifyPageUserBefore(user._id);
        },
      )
      .waitForFirstXHR(
        API_ROUTES.user.list,
        WAIT_FOR_FIRST_XHR_TIME,
        () => usersPage.clickClearSearchButton(), ({ responseData }) => {
          const { data } = JSON.parse(responseData);

          browser.assert.ok(data.some(item => item._id !== user._id));
        },
      );
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

    const userSelector = users.create._id;

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

    browser.completed.user.delete(createUser, ({ responseData }) => {
      users.create = {
        ...users.create,
        deleteData: responseData,
      };
    });

    browser.completed.user.delete(editUser, ({ responseData }) => {
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

    users.mass.map((user, index) => browser.completed.user.create(user, ({ userResponseData }) => {
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
