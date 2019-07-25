// http://nightwatchjs.org/guide#usage

module.exports.command = function deleteView(
  view = {
    title: 'test-title',
    tags: 'test-tags',
  },
  callback = result => result,
) {
  const {
    title = 'test-title',
    tags = 'test-tags',
  } = view;

  const viewPath = `${process.env.VUE_DEV_SERVER_URL}view/`;

  const topBar = this.page.layout.topBar();
  const confirmation = this.page.modals.confirmation();
  const leftSideBar = this.page.layout.leftSideBar();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const modalViewCreate = this.page.modals.view.create();

  this.url(viewPath);

  topBar.clickUserDropdown()
    .clickUserProfileButton();

  createUser.verifyModalOpened()
    .selectNavigationType(1)
    .clickSubmitButton()
    .verifyModalClosed();

  groupsSideBar.clickGroupsSideBarButton();

  leftSideBar.verifySettingsWrapperBefore()
    .clickSettingsViewButton()
    .verifyControlsWrapperBefore()
    .clickEditViewButton()
    .defaultPause();

  groupsSideBar.clickPanelHeader(tags)
    .verifyPanelBody(tags)
    .clickEditViewButton(title)
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .clickViewDeleteButton();

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    '/api/v2/views/',
    1000,
    () => {
      confirmation.clickConfirmButton();
    },
    (xhr) => {
      callback({ xhr });
    },
  );

  confirmation.verifyModalClosed();

  return this;
};
