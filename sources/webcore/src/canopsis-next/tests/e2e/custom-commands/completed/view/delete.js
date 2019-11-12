// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

module.exports.command = function deleteView(groupId, viewId, callback = () => {}) {
  const topBar = this.page.layout.topBar();
  const confirmation = this.page.modals.common.confirmation();
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

  groupsSideBar.browseGroupById(groupId)
    .verifyPanelBody(groupId)
    .clickEditViewButton(viewId)
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .clickViewDeleteButton();

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    `${API_ROUTES.view}/${viewId}`,
    WAIT_FOR_FIRST_XHR_TIME,
    () => confirmation.clickSubmitButton(),
    ({ responseData }) => {
      confirmation.verifyModalClosed();

      callback(JSON.parse(responseData));
    },
  );

  return this;
};
