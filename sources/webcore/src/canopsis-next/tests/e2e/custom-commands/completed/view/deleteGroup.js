// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function deleteGroup(groupId, callback = () => {}) {
  const topBar = this.page.layout.topBar();
  const navigation = this.page.layout.navigation();
  const createUser = this.page.modals.admin.createUser();
  const confirmation = this.page.modals.common.confirmation();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const modalCreateGroup = this.page.modals.view.createGroup();

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

  groupsSideBar.clickEditGroupButton(groupId)
    .defaultPause();

  modalCreateGroup.verifyModalOpened()
    .clickDeleteButton();

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    `${API_ROUTES.viewGroup}/${groupId}`,
    5000,
    () => confirmation.clickSubmitButton(),
    ({ responseData }) => {
      confirmation.verifyModalClosed();

      callback(JSON.parse(responseData));
    },
  );
};
