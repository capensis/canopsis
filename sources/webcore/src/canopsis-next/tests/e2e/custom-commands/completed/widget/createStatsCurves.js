// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');

module.exports.command = function createStatsCurves(fields, callback = () => {}) {
  const statsCurvesWidget = this.page.widget.statsCurves();

  this.completed.widget.setCommonFields(fields);

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => statsCurvesWidget.clickSubmitStatsCurves(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
