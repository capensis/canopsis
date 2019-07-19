// http://nightwatchjs.org/guide#usage

const { NAVIGATION: { SIDEBAR } } = require('../../constants');

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
    browser.end(done);
  },

  'Open current user modal': (browser) => {
    const topBar = browser.page.layout.topBar();
    const popup = browser.page.layout.popup();
    const createUserModal = browser.page.modals.admin.createUser();

    popup.clickOnEveryPopupsCloseIcons();

    topBar.clickUserDropdown()
      .clickUserProfileButton();

    createUserModal.verifyModalOpened();
  },

  'Select current user default view': (browser) => {
    const createUserModal = browser.page.modals.admin.createUser();
    const selectViewModal = browser.page.modals.view.selectView();

    createUserModal.clickSelectDefaultViewButton();

    selectViewModal.verifyModalOpened()
      .browseGroupById('05b2e049-b3c4-4c5b-94a5-6e7ff142b28c') // TODO: use from some constants file when we will use fixtures
      .browseViewById('875df4c2-027b-4549-8add-e20ed7ff7d4f')
      .verifyModalClosed();

    createUserModal.clickSubmitButton()
      .verifyModalClosed();
  },

  'Check default view': (browser) => {
    browser.url(process.env.VUE_DEV_SERVER_URL)
      .page.view()
      .verifyPageElementsBeforeById('875df4c2-027b-4549-8add-e20ed7ff7d4f'); // TODO: use from some constants file when we will use fixtures
  },

  'Add view with some name from constants': (browser) => {
    const { text, create: { prefix } } = SIDEBAR;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description, tags,
    } = TEMPORARY_DATA[prefix];

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
    const { text, copy: { prefix }, create } = SIDEBAR;

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
    const { text, edit: { prefix }, create } = SIDEBAR;

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
    const { create, edit, copy } = SIDEBAR;

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
      .verifyModalClosed();
  },
};
