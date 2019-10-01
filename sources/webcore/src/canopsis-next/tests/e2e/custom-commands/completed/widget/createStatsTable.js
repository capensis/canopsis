// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');

module.exports.command = function createStatsTable({
  parameters,
  ...fields
}, callback = () => {}) {
  const statsTableWidget = this.page.widget.statsTable();

  this.completed.widget.setCommonFields({
    parameters: {
      advanced: true,
      ...parameters,
    },
    ...fields,
  });

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => statsTableWidget.clickSubmitStatsTable(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
