// http://nightwatchjs.org/guide#usage
const uid = require('uid');

const { createAdminUser, removeUser } = require('../../helpers/api');
const { generateTemporaryView } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.views = [];

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

  'Add view with some name from constants': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.views.push(view);
    });
  },

  'Checking view copy with name from constants': (browser) => {
    const [createdView] = browser.globals.views;

    const { title, name, enabled } = generateTemporaryView();

    browser.completed.view.copy(createdView.group_id, createdView._id, { title, name, enabled }, (view) => {
      browser.globals.views.push(view);
    });
  },

  'Editing test view with name from constants': (browser) => {
    const navigation = browser.page.layout.navigation();
    const groupsSideBar = browser.page.layout.groupsSideBar();
    const modalCreateGroup = browser.page.modals.view.createGroup();

    const [createdView] = browser.globals.views;
    const groupName = `group-${uid()}`;

    navigation.verifySettingsWrapperBefore()
      .clickSettingsViewButton()
      .verifyControlsWrapperBefore()
      .clickEditModeButton()
      .defaultPause();

    groupsSideBar.clickEditGroupButton(createdView.group_id)
      .defaultPause();

    modalCreateGroup.verifyModalOpened()
      .clearGroupName()
      .setGroupName(groupName)
      .clickSubmitButton()
      .verifyModalClosed();

    navigation.clickEditModeButton()
      .clickSettingsViewButton();

    browser.completed.view.edit(createdView.group_id, createdView._id, generateTemporaryView(), (view) => {
      /**
       * TODO: put group removing when it will be ready
       */
      browser.globals.views[0] = view;
    });
  },

  'Deleting all test items view with name from constants': (browser) => {
    browser.globals.views.forEach(view => browser.completed.view.delete(view.group_id, view._id));
    browser.globals.views.forEach(view => browser.completed.view.deleteGroup(view.group_id));
  },
};
