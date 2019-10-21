// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

module.exports.command = function createStatsNumber({
  parameters: {
    statSelector,
    displayMode,
    sortOrder,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const statsNumberWidget = this.page.widget.statsNumber();
  const statsDisplayModeModal = this.page.modals.stats.statsDisplayMode();
  const colorPickerModal = this.page.modals.common.colorPicker();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters: {
      advanced: true,
      ...parameters,
    },
  });

  if (displayMode) {
    statsNumberWidget.clickDisplayModeEditButton();

    statsDisplayModeModal
      .verifyModalOpened()
      .selectDisplayModeType(displayMode.type);

    Object.entries(displayMode.parameters).forEach(([level, { value, color }]) => {
      statsDisplayModeModal
        .clickDisplayModeParameterValue(level)
        .clearDisplayModeParameterValue(level)
        .setDisplayModeParameterValue(level, value);

      statsDisplayModeModal.clickDisplayModeParameterColor(level);

      colorPickerModal
        .verifyModalOpened()
        .clickColorField()
        .clearColorField()
        .setColorField(color)
        .clickSubmitButton()
        .verifyModalClosed();
    });

    statsDisplayModeModal
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (sortOrder) {
    statsNumberWidget
      .clickSortOrder()
      .selectSortOrder(sortOrder);
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    WAIT_FOR_FIRST_XHR_TIME,
    () => statsNumberWidget.clickSubmitStatsNumber(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
