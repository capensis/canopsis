// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');

module.exports.command = function createStatsHistogram(fields, callback = () => {}) {
  const statsHistogramWidget = this.page.widget.statsHistogram();

  this.completed.widget.setCommonFields(fields);

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => statsHistogramWidget.clickSubmitStatsHistogram(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
