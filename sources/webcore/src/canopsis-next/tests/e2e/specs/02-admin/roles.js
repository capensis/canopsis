const createRole = (browser, {
  name, description, groupId, viewId,
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

  createRoleModal.clickSubmitButton()
    .verifyModalClosed();
};

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Create new role with data from constants': (browser) => {
    const role = {
      name: 'Test role',
      description: 'Test role description',
      groupId: '05b2e049-b3c4-4c5b-94a5-6e7ff142b28c',
      viewId: '875df4c2-027b-4549-8add-e20ed7ff7d4f',
    };

    browser.page.admin.roles()
      .navigate()
      .verifyPageElementsBefore();

    createRole(browser, role);
  },
};
