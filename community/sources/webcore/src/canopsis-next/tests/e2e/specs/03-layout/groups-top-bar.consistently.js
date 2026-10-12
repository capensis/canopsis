// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../src/config');
const { NAVIGATION_TYPES, WAIT_FOR_FIRST_XHR_TIME } = require('../../constants');
const { generateTemporaryView } = require('../../helpers/entities');
const { createAdminUser, removeUser } = require('../../helpers/api');

module.exports = {
  async before(browser, done) {
    browser.globals.views = {};

    const { data } = await createAdminUser();

    browser.globals.credentials = {
      password: data.password,
      username: data._id,
    };

    await browser.maximizeWindow()
      .completed.login(browser.globals.credentials.username, browser.globals.credentials.password);

    done();
  },

  async after(browser, done) {
    browser.completed.logout()
      .end();

    await removeUser(browser.globals.credentials.username);

    delete browser.globals.credentials;
    delete browser.globals.views;

    done();
  },

  'Add view with name from constants': (browser) => {
    const { views } = browser.globals;

    views.create = generateTemporaryView('create');

    const {
      name, title, description, group,
    } = views.create;

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectNavigationType(NAVIGATION_TYPES.topBar)
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
      WAIT_FOR_FIRST_XHR_TIME,
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
      WAIT_FOR_FIRST_XHR_TIME,
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
      WAIT_FOR_FIRST_XHR_TIME,
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
    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickEditGroupButton(views.create.group_id)
      .defaultPause();
    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clickDeleteButton();
    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickDropdownButton(views.edit.group_id)
      .verifyDropdownZone(views.edit.group_id)
      .clickEditViewButton(views.edit._id)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clickViewDeleteButton();

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickEditGroupButton(views.edit.group_id)
      .defaultPause();

    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clickDeleteButton();

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Reset user data': (browser) => {
    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectNavigationType(NAVIGATION_TYPES.sideBar)
      .clickSubmitButton()
      .verifyModalClosed();
  },
};
