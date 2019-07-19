// http://nightwatchjs.org/guide#usage

const faker = require('faker');

const TEMPORARY_DATA = {};

const createTemporaryObject = () => ({
  name: faker.lorem.word(),
  title: faker.lorem.words(),
  description: faker.lorem.sentences(),
  tags: faker.lorem.slug(),
});

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

  'Create test view': (browser) => {
    const topBar = browser.page.layout.topBar();
    const createUserModal = browser.page.modals.admin.createUser();

    topBar.clickUserDropdown()
      .clickUserProfileButton();

    createUserModal.verifyModalOpened()
      .selectNavigationType(1)
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Add view with some name from constants': (browser) => {
    TEMPORARY_DATA.create = createTemporaryObject();

    const {
      name, title, description, tags,
    } = TEMPORARY_DATA.create;

    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton();

    browser.page.layout.leftSideBar()
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
      .setViewGroupTags(tags)
      .setViewGroupIds(tags)
      .clickViewSubmitButton()
      .verifyModalClosed();
  },
};
