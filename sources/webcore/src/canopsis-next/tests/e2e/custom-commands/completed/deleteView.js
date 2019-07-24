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

  const confirmation = this.page.modals.confirmation();
  const leftSideBar = this.page.layout.leftSideBar();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const modalViewCreate = this.page.modals.view.create();

  leftSideBar.verifyControlsWrapperBefore()
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
