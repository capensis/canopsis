// http://nightwatchjs.org/guide#usage

const { NAVIGATION: { groups } } = require('../../constants');

const TEMPORARY_DATA = {};

const createTemporaryObject = ({ prefix, text, index }) => {
  const i = typeof index === 'number' ? `-${index}` : '';
  const r = Math.random().toString(36).substring(7);
  return {
    name: `${prefix}-${text}-name${i}-${r}`,
    title: `${prefix}-${text}-title${i}-${r}`,
    description: `${prefix}-${text}-description${i}-${r}`,
    tags: `${prefix}-${text}-tags${i}-${r}`,
  };
};

module.exports = {
  async before(browser, done) {
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

    browser.end(done);
  },

  'Add view with name from constants': (browser) => {
    const { text, create: { prefix } } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description, tags,
    } = TEMPORARY_DATA[prefix];

    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectNavigationType(2)
      .clickSubmitButton()
      .verifyModalClosed()
      .api.pause(5000);

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

  'Checking view copy with name from constants': (browser) => {
    const { text, copy: { prefix }, create } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description,
    } = TEMPORARY_DATA[prefix];

    browser.page.layout.leftSideBar()
      .verifyControlsWrapperBefore()
      .clickEditViewButton()
      .defaultPause();

    browser.page.layout.topBar()
      .clickDropdownButton(TEMPORARY_DATA[create.prefix].tags)
      .verifyDropdownZone(TEMPORARY_DATA[create.prefix].tags)
      .clickCopyViewButton(TEMPORARY_DATA[create.prefix].title)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .setViewName(name)
      .setViewTitle(title)
      .clearViewDescription()
      .setViewDescription(description)
      .clickViewSubmitButton()
      .verifyModalClosed();
  },

  'Editing test view with name from constants': (browser) => {
    const { text, edit: { prefix }, create } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description, tags,
    } = TEMPORARY_DATA[prefix];

    const r = Math.random().toString(36).substring(7);

    browser.page.layout.topBar()
      .clickEditGroupButton(TEMPORARY_DATA[create.prefix].tags)
      .defaultPause();

    TEMPORARY_DATA[create.prefix].tags = `${create.prefix}-${text}-tags-${r}`;

    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clearGroupName()
      .setGroupName(TEMPORARY_DATA[create.prefix].tags)
      .clickSubmitButton()
      .verifyModalClosed();


    browser.page.layout.topBar()
      .clickDropdownButton(TEMPORARY_DATA[create.prefix].tags)
      .verifyDropdownZone(TEMPORARY_DATA[create.prefix].tags)
      .clickEditViewButton(TEMPORARY_DATA[create.prefix].title)
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
      .setViewGroupTags(tags)
      .clearViewGroupIds()
      .setViewGroupIds(tags)
      .clickViewSubmitButton()
      .verifyModalClosed();
  },

  'Deleting all test items view with name from constants': (browser) => {
    const { create, edit, copy } = groups;

    browser.page.layout.topBar()
      .clickDropdownButton(TEMPORARY_DATA[create.prefix].tags)
      .verifyDropdownZone(TEMPORARY_DATA[create.prefix].tags)
      .clickEditViewButton(TEMPORARY_DATA[copy.prefix].title)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clickViewDeleteButton();
    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();

    browser.page.layout.topBar()
      .clickDropdownButton(TEMPORARY_DATA[edit.prefix].tags)
      .verifyDropdownZone(TEMPORARY_DATA[edit.prefix].tags)
      .clickEditViewButton(TEMPORARY_DATA[edit.prefix].title)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clickViewDeleteButton();
    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed()
      .api.pause(5000);
  },
};
