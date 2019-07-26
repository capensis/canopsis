// http://nightwatchjs.org/guide#usage
const { VIEW } = require('../../constants');
const { API_ROUTES } = require('../../../../src/config');

module.exports.command = function deleteView(
  view = VIEW,
  callback = result => result,
) {
  const {
    title = VIEW.title,
    group = VIEW.group,
  } = view;

  const topBar = this.page.layout.topBar();
  const confirmation = this.page.modals.confirmation();
  const navigation = this.page.layout.navigation();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const modalViewCreate = this.page.modals.view.create();

  this.refresh();

  topBar.clickUserDropdown()
    .clickUserProfileButton();

  createUser.verifyModalOpened()
    .selectNavigationType(1)
    .clickSubmitButton()
    .verifyModalClosed();

  groupsSideBar.clickGroupsSideBarButton();

  navigation.verifySettingsWrapperBefore()
    .clickSettingsViewButton()
    .verifyControlsWrapperBefore()
    .clickEditViewButton()
    .defaultPause();

  groupsSideBar.clickPanelHeader(group)
    .verifyPanelBody(group)
    .clickEditViewButton(title)
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .clickViewDeleteButton();

  confirmation.verifyModalOpened();

  this.waitForFirstXHR(
    API_ROUTES.view,
    1000,
    () => confirmation.clickConfirmButton(),
    ({ responseData }) => {
      callback(JSON.parse(responseData));
    },
  );

  confirmation.verifyModalClosed();

  return this;
};
