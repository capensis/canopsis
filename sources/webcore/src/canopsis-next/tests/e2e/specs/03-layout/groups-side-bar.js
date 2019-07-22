// http://nightwatchjs.org/guide#usage

const { NAVIGATION: { LEFT_SIDEBAR } } = require('../../constants');

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

  'Browse view by id': (browser) => {
    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .browseGroupById('05b2e049-b3c4-4c5b-94a5-6e7ff142b28c') // TODO: use from some constants file when we will use fixtures
      .browseViewById('da7ac9b9-db1c-4435-a1f2-edb4d6be4db8')
      .defaultPause(); // TODO: put verification
  },

  'Add view with some name from constants': (browser) => {
    const { text, create: { prefix } } = LEFT_SIDEBAR;

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
    const { text, copy: { prefix }, create } = LEFT_SIDEBAR;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description,
    } = TEMPORARY_DATA[prefix];

    browser.page.layout.leftSideBar()
      .verifyControlsWrapperBefore()
      .clickEditViewButton()
      .defaultPause();

    browser.page.layout.groupsSideBar()
      .clickPanelHeader(TEMPORARY_DATA[create.prefix].tags)
      .verifyPanelBody(TEMPORARY_DATA[create.prefix].tags)
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
    const { text, edit: { prefix }, create } = LEFT_SIDEBAR;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description, tags,
    } = TEMPORARY_DATA[prefix];

    const r = Math.random().toString(36).substring(7);

    browser.page.layout.groupsSideBar()
      .clickEditGroupButton(TEMPORARY_DATA[create.prefix].tags)
      .defaultPause();

    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clearGroupName()
      .setGroupName(`${create.prefix}-${text}-tags-${r}`)
      .clickSubmitButton()
      .verifyModalClosed();


    browser.page.layout.groupsSideBar()
      .verifyPanelBody(`${create.prefix}-${text}-tags-${r}`)
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
    const { edit, copy } = LEFT_SIDEBAR;

    browser.page.layout.groupsSideBar()
      .clickEditViewButton(TEMPORARY_DATA[copy.prefix].title)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clickViewDeleteButton();
    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();

    browser.page.layout.groupsSideBar()
      .clickPanelHeader(TEMPORARY_DATA[edit.prefix].tags)
      .verifyPanelBody(TEMPORARY_DATA[edit.prefix].tags)
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
