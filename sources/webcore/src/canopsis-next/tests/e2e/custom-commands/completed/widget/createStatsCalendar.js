// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createStatsCalendar({
  parameters: {
    typeOfEntities,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const statsCalendarWidget = this.page.widget.statsCalendar();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters: {
      alarmsList: true,
      advanced: true,
      ...parameters,
    },
  });

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
