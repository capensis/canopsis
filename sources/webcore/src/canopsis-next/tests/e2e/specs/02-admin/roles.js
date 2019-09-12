const { API_ROUTES } = require('../../../../src/config');
const { generateTemporaryView, generateTemporaryRole } = require('../../helpers/entities');

const createRole = (browser, {
  name,
  description,
  groupId,
  viewId,
}) => {
  const rolesPage = browser.page.admin.roles();
  const createRoleModal = browser.page.modals.admin.createRole();
  const selectViewModal = browser.page.modals.view.selectView();

  rolesPage.clickAddButton();

  createRoleModal.verifyModalOpened()
    .setName(name)
    .setDescription(description)
    .clickSelectDefaultViewButton();

  selectViewModal.verifyModalOpened()
    .browseGroupById(groupId)
    .browseViewById(viewId)
    .verifyModalClosed();

  browser.waitForFirstXHR(
    API_ROUTES.role.create,
    1000,
    () => createRoleModal.clickSubmitButton(),
    ({ responseData }) => {
      const response = JSON.parse(responseData);

      browser.assert.equal(response.total, 1);

      browser.globals.roles.push(response.data[0]);

      createRoleModal.verifyModalClosed();
    },
  );
};

module.exports = {
  async before(browser, done) {
    browser.globals.roles = [];
    browser.globals.defaultViewData = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    delete browser.globals.defaultViewData;
    delete browser.globals.roles;

    browser.completed.logout()
      .end(done);
  },

  'Create test view': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.defaultViewData = {
        viewId: view._id,
        groupId: view.group_id,
      };
    });
  },

  'Create new role with data from constants': (browser) => {
    const generatedRole = {
      ...generateTemporaryRole(),
      ...browser.globals.defaultViewData,
    };

    browser.page.admin.roles()
      .navigate()
      .verifyPageElementsBefore();

    createRole(browser, generatedRole);
  },


  'Check searching': (browser) => {
    const [role] = browser.globals.roles;
    const rolesPage = browser.page.admin.roles();

    rolesPage.setSearchingText(role._id)
      .waitForFirstXHR(API_ROUTES.role.list, 5000, () => rolesPage.clickSubmitSearchButton(), ({ responseData }) => {
        const { data } = JSON.parse(responseData);

        browser.assert.ok(data.every(item => item._id === role._id));
        browser.assert.elementsCount(rolesPage.elements.dataTableUserItem.selector, 1);

        rolesPage.verifyPageRoleBefore(role._id);
      })
      .waitForFirstXHR(API_ROUTES.role.list, 5000, () => rolesPage.clickClearSearchButton(), ({ responseData }) => {
        const { data } = JSON.parse(responseData);

        browser.assert.ok(data.some(item => item._id !== role._id));
      });
  },

  'Pagination on data-table': (browser) => {
    const rolesPage = browser.page.admin.roles();

    rolesPage.clickNextButton()
      .defaultPause();

    rolesPage.clickPrevButton()
      .defaultPause();

    rolesPage.selectRange(5)
      .defaultPause();
  },

  'Edit created role by data from constants': (browser) => {
    const rolesPage = browser.page.admin.roles();
    const createRoleModal = browser.page.modals.admin.createRole();
    const { _id: roleId } = browser.globals.roles[0];

    rolesPage.verifyPageRoleBefore(roleId)
      .clickEditButton(roleId);

    createRoleModal.verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Delete created role': (browser) => {
    const rolesPage = browser.page.admin.roles();
    const confirmationModal = browser.page.modals.common.confirmation();
    const { _id: roleId } = browser.globals.roles.shift();

    rolesPage.verifyPageRoleBefore(roleId)
      .clickDeleteButton(roleId);

    confirmationModal.verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Create several new roles with data from constants': (browser) => {
    const roles = [];

    for (let index = 0; index < 3; index += 1) {
      roles.push({
        ...generateTemporaryRole(),
        ...browser.globals.defaultViewData,
      });
    }

    browser.page.admin.roles()
      .navigate()
      .verifyPageElementsBefore();

    roles.forEach((role) => {
      createRole(browser, role);
    });
  },

  'Mass delete created roles': (browser) => {
    const rolesPage = browser.page.admin.roles();
    const confirmationModal = browser.page.modals.common.confirmation();
    const { roles } = browser.globals;

    rolesPage.selectRange(5)
      .defaultPause();

    roles.forEach((role) => {
      rolesPage.verifyPageRoleBefore(role._id)
        .clickOptionCheckbox(role._id)
        .defaultPause();
    });

    rolesPage.verifyMassDeleteButton()
      .clickMassDeleteButton();

    confirmationModal.verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Refresh button': (browser) => {
    const rolesPage = browser.page.admin.roles();

    browser.completed.refreshPage(API_ROUTES.role.list, () => rolesPage.clickRefreshButton());
  },

  'Delete test view': (browser) => {
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.completed.view.delete(groupId, viewId);
    browser.completed.view.deleteGroup(groupId);
  },
};
