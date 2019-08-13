// http://nightwatchjs.org/guide#usage

const qs = require('qs');
const {
  // generateTemporaryUser,
  // generateTemporaryRight,
  generateTemporaryRole,
} = require('../../helpers/entities');

const { API_ROUTES } = require('../../../../src/config');

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  async after(browser, done) {
    delete browser.globals.temporary;

    await browser.completed.logout()
      .end(done);
  },

  'Check tabs': (browser) => {
    browser.page.admin.right()
      .navigate();

    browser.page.admin.right()
      .clickTab('view')
      .api.pause(1000);

    browser.page.admin.right()
      .clickTab('technical')
      .api.pause(1000);

    browser.page.admin.right()
      .clickTab('business')
      .api.pause(1000);

    browser.refresh();
  },

  // 'Create new right': (browser) => {
  //   const { temporary } = browser.globals;

  //   temporary.right = generateTemporaryRight();

  //   const {
  //     ID,
  //     description,
  //   } = temporary.right;

  //   const adminCreateRight = browser.page.modals.admin.createRight();

  //   browser.page.admin.right()
  //     .clickAddButton()
  //     .clickCreateRight();

  //   adminCreateRight.verifyModalOpened()
  //     .setRightID(ID)
  //     .setRightDescription(description);

  //   browser.waitForFirstXHR(
  //     API_ROUTES.action,
  //     1000,
  //     () => adminCreateRight.clickSubmitButton(),
  //     ({ responseData, requestData }) => {
  //       const requestParsedData = qs.parse(requestData);

  //       return temporary.right = {
  //         ...temporary.right,
  //         ...JSON.parse(requestParsedData),
  //         ...JSON.parse(responseData),
  //       };
  //     },
  //   );
  //   adminCreateRight.verifyModalClosed();
  // },

  'Create new user role (put special rights)': (browser) => {
    const { temporary } = browser.globals;

    temporary.role = generateTemporaryRole();

    const {
      name,
      description,
    } = temporary.role;

    const adminCreateRole = browser.page.modals.admin.createRole();

    browser.page.admin.right()
      .clickAddButton()
      .clickCreateRole();

    adminCreateRole.verifyModalOpened()
      .setRoleName(name)
      .setRoleDescription(description);

    browser.waitForFirstXHR(
      API_ROUTES.role.create,
      5000,
      () => adminCreateRole.clickSubmitButton(),
      ({ responseData, requestData }) => {
        const requestParsedData = qs.parse(requestData);
        return temporary.role = {
          ...temporary.role,
          ...JSON.parse(requestParsedData.role),
          ...JSON.parse(responseData),
        };
      },
    );
    adminCreateRole.verifyModalClosed();
  },

  // 'Create user with this new role': (browser) => {
  //   const { temporary } = browser.globals;

  //   temporary.user = generateTemporaryUser('user');

  //   const {
  //     username,
  //     firstname,
  //     lastname,
  //     email,
  //     password,
  //   } = temporary.user;

  //   const adminCreateUser = browser.page.modals.admin.createUser();


  //   browser.page.admin.right()
  //     .clickCreateUser();

  //   adminCreateUser.verifyModalOpened()
  //     .setUsername(username)
  //     .setFirstName(firstname)
  //     .setLastName(lastname)
  //     .setEmail(email)
  //     .setPassword(password)
  //     .selectLastRole()
  //     .selectLanguage()
  //     .selectNavigationType();

  //   browser.waitForFirstXHR(
  //     API_ROUTES.user.create,
  //     5000,
  //     () => adminCreateUser.clickSubmitButton(),
  //     ({ responseData, requestData }) => {
  //       const requestParsedData = qs.parse(requestData);

  //       return temporary.user = {
  //         ...temporary.user,
  //         ...JSON.parse(requestParsedData.user),
  //         ...JSON.parse(responseData),
  //       };
  //     },
  //   );
  //   adminCreateUser.verifyModalClosed();
  // },


  // 'Login by new role user credentials': (browser) => {
  //   const { username, password } = browser.globals.user;
  //   browser.completed.logout()
  //     .maximizeWindow()
  //     .completed.loginDisabledUser(username, password);

  //   browser.api.pause(10000);
  // },

  'Change right for new role': (browser) => {
    const { role } = browser.globals.temporary;

    browser.page.admin.right()
      .clickCheckbox(role.name, 'listalarm_ack', 1);

    browser.page.admin.right()
      .clickTab('view');

    browser.page.admin.right()
      .clickCheckbox(role.name, '875df4c2-027b-4549-8add-e20ed7ff7d4f', 1)
      .clickCheckbox(role.name, '875df4c2-027b-4549-8add-e20ed7ff7d4f', 2)
      .clickCheckbox(role.name, '875df4c2-027b-4549-8add-e20ed7ff7d4f', 3);

    browser.page.admin.right()
      .clickTab('technical');

    browser.page.admin.right()
      .clickCheckbox(role.name, 'models_userview', 1)
      .clickCheckbox(role.name, 'models_userview', 2)
      .clickCheckbox(role.name, 'models_userview', 3)
      .clickCheckbox(role.name, 'models_userview', 4);

    browser.page.admin.right().clickSubmitRightButton();

    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();
  },

};
