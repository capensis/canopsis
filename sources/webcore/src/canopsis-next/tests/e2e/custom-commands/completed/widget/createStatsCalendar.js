// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createStatsCalendar({
  parameters: {
    criticityLevels,
    colorsSelector,
    considerPbehaviors,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const statsCalendarWidget = this.page.widget.statsCalendar();
  const colorPickerModal = this.page.modals.common.colorPicker();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters: {
      alarmsList: true,
      advanced: true,
      ...parameters,
    },
  });

  if (criticityLevels) {
    statsCalendarWidget
      .clickCriticityLevels()
      .clickCriticityLevelsMinor()
      .clearCriticityLevelsMinor()
      .setCriticityLevelsMinor(criticityLevels.minor)
      .clickCriticityLevelsMajor()
      .clearCriticityLevelsMajor()
      .setCriticityLevelsMajor(criticityLevels.major)
      .clickCriticityLevelsCritical()
      .clearCriticityLevelsCritical()
      .setCriticityLevelsCritical(criticityLevels.critical);
  }

  if (colorsSelector) {
    statsCalendarWidget.clickColorSelector();

    Object.entries(colorsSelector).forEach(([level, color]) => {
      statsCalendarWidget.clickColorPickerButton(level);

      colorPickerModal
        .verifyModalOpened()
        .clickColorField()
        .clearColorField()
        .setColorField(color)
        .clickSubmitButton()
        .verifyModalClosed();
    });
  }

  if (typeof considerPbehaviors === 'boolean') {
    statsCalendarWidget.setConsiderPbehaviors(considerPbehaviors);
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => statsCalendarWidget.clickSubmitStatsCalendar(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
