// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createAlarmsList({
  parameters: {
    ack,
    moreInfos,
    enableHtml = false,
    infoPopups,
    newColumnNames,
    editColumnNames,
    moveColumnNames,
    deleteColumnNames,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const common = this.page.widget.common();
  const textEditor = this.page.modals.common.textEditor();
  const addInfoPopup = this.page.modals.common.addInfoPopup();
  const infoPopupModal = this.page.modals.common.infoPopupSetting();
  // const createFilter = this.page.modals.common.createFilter();
  const alarms = this.page.widget.alarms();

  this.completed.widget.setCommonFields({ ...fields, parameters });

  if (newColumnNames || editColumnNames || moveColumnNames || deleteColumnNames) {
    common.clickColumnNames();
  }

  if (newColumnNames) {
    newColumnNames.forEach(({ index, data }) => {
      common
        .clickAddColumnName()
        .editColumnName(index, data);
    });
  }

  if (editColumnNames) {
    editColumnNames.forEach(({ index, data }) => {
      common.editColumnName(index, data);
    });
  }

  if (moveColumnNames) {
    moveColumnNames.forEach(({ index, up, down }) => {
      if (up) {
        common.clickColumnNameUpWard(index);
      }

      if (down) {
        common.clickColumnNameDownWard(index);
      }
    });
  }

  if (deleteColumnNames) {
    deleteColumnNames.forEach((index) => {
      common.clickDeleteColumnName(index);
    });
  }

  if (moreInfos) {
    common.clickCreateMoreInfos();

    textEditor.verifyModalOpened()
      .clickField()
      .setField(moreInfos)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (enableHtml) {
    alarms.setEnableHtml(enableHtml);
  }

  if (ack) {
    alarms.clickAckGroup()
      .setIsAckNoteRequired(ack.isAckNoteRequired)
      .setIsMultiAckEnabled(ack.isMultiAckEnabled);

    if (ack.fastAckOutput) {
      alarms.clickFastAckOutput()
        .setFastAckOutputSwitch(ack.fastAckOutput.enabled);
    }

    if (ack.fastAckOutput.enabled) {
      alarms.clickFastAckOutputText()
        .clearFastAckOutputText()
        .setFastAckOutputText(ack.fastAckOutput.output);
    }
  }

  if (infoPopups) {
    common.clickInfoPopup();

    infoPopupModal.verifyModalOpened();

    infoPopups.forEach(({ field, template }) => {
      infoPopupModal.clickAddPopup();

      addInfoPopup.verifyModalOpened()
        .selectSelectedColumn(field)
        .setTemplate(template)
        .clickSubmitButton()
        .verifyModalClosed();
    });

    infoPopupModal.clickSubmitButton()
      .verifyModalClosed();
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => alarms.clickSubmitAlarms(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
