// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../src/config');

module.exports.command = function deleteView(groupId, viewId, callback = () => {}) {
  const topBar = this.page.layout.topBar();
  const confirmation = this.page.modals.confirmation();
  const navigation = this.page.layout.navigation();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const modalViewCreate = this.page.modals.view.create();

  groupsSideBar.groupsSideBarButtonElement(({ status }) => {
    if (status === -1) {
      topBar.clickUserDropdown()
        .clickUserProfileButton();

      createUser.verifyModalOpened()
        .selectNavigationType(1)
        .clickSubmitButton()
        .verifyModalClosed();
    }

    groupsSideBar.clickGroupsSideBarButton();
  });

  navigation.verifySettingsWrapperBefore()
    .clickSettingsViewButton()
    .verifyControlsWrapperBefore()
    .clickEditModeButton()
    .defaultPause();

  groupsSideBar.clickPanelHeader(groupId)
    .verifyPanelBody(groupId)
    .clickEditViewButton(viewId)
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .clickViewDeleteButton();

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    new RegExp(`${API_ROUTES.view}$`),
    5000,
    () => confirmation.clickConfirmButton(),
    ({ responseData }) => callback(JSON.parse(responseData)),
  );

  confirmation.verifyModalClosed();

  navigation.clickEditModeButton()
    .clickSettingsViewButton();

  return this;
};
