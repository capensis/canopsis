// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../src/config');
const { generateTemporaryView } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.views = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectNavigationType(1)
      .clickSubmitButton()
      .verifyModalClosed();

    delete browser.globals.views;

    browser.end(done);
  },

  'Add view with name from constants': (browser) => {
    const { views } = browser.globals;

    views.create = generateTemporaryView('create');

    const {
      name, title, description, group,
    } = views.create;

    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectNavigationType(2)
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.navigation()
      .verifySettingsWrapperBefore()
      .clickSettingsViewButton()
      .verifyControlsWrapperBefore()
      .clickAddViewButton()
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .setViewName(name)
      .setViewTitle(title)
      .setViewDescription(description)
      .clickViewEnabled()
      .setViewGroupTags(group)
      .setViewGroupId(group);

    browser.waitForFirstXHR(
      new RegExp(`${API_ROUTES.view}$`),
      5000,
      () => browser.page.modals.view.create()
        .clickViewSubmitButton(),
      ({ responseData, requestData }) => views.create = {
        ...views.create,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );


    browser.page.modals.view.create()
      .verifyModalClosed();
  },

  'Checking view copy with name from constants': (browser) => {
    const { views } = browser.globals;

    views.copy = generateTemporaryView('copy');

    const {
      name, title, description,
    } = views.copy;

    browser.page.layout.navigation()
      .verifyControlsWrapperBefore()
      .clickEditModeButton()
      .defaultPause();

    browser.page.layout.topBar()
      .clickDropdownButton(views.create.group_id)
      .verifyDropdownZone(views.create.group_id)
      .clickCopyViewButton(views.create._id)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .setViewName(name)
      .setViewTitle(title)
      .clearViewDescription()
      .setViewDescription(description);

    browser.waitForFirstXHR(
      new RegExp(`${API_ROUTES.view}$`),
      5000,
      () => browser.page.modals.view.create()
        .clickViewSubmitButton(),
      ({ responseData, requestData }) => views.copy = {
        ...views.copy,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    browser.page.modals.view.create()
      .verifyModalClosed();
  },

  'Editing test view with name from constants': (browser) => {
    const { views } = browser.globals;

    views.edit = generateTemporaryView('edit');

    const {
      name, title, description, group,
    } = views.edit;

    browser.page.layout.topBar()
      .clickEditGroupButton(views.create.group_id)
      .defaultPause();

    views.create.group = generateTemporaryView('create').group;

    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clearGroupName()
      .setGroupName(views.create.group)
      .clickSubmitButton()
      .verifyModalClosed();


    browser.page.layout.topBar()
      .clickDropdownButton(views.create.group_id)
      .verifyDropdownZone(views.create.group_id)
      .clickEditViewButton(views.create._id)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clearViewName()
      .setViewName(name)
      .clearViewTitle()
      .setViewTitle(title)
      .clearViewDescription()
      .setViewDescription(description)
      .clickViewEnabled()
      .clearViewGroupTags()
      .setViewGroupTags(group)
      .clearViewGroupId()
      .setViewGroupId(group);


    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.create._id}`,
      5000,
      () => browser.page.modals.view.create()
        .clickViewSubmitButton(),
      ({ responseData, requestData }) => views.edit = {
        ...views.edit,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    browser.page.modals.view.create()
      .verifyModalClosed();
  },

  'Deleting all test items view with name from constants': (browser) => {
    const { views } = browser.globals;

    browser.page.layout.topBar()
      .clickDropdownButton(views.create.group_id)
      .verifyDropdownZone(views.create.group_id)
      .clickEditViewButton(views.copy._id)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clickViewDeleteButton();
    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();

    browser.page.layout.topBar()
      .clickDropdownButton(views.edit.group_id)
      .verifyDropdownZone(views.edit.group_id)
      .clickEditViewButton(views.edit._id)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clickViewDeleteButton();
    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },
};
