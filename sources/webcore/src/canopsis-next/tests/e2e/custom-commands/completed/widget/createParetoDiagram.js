// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');

module.exports.command = function createParetoDiagram(fields, callback = () => {}) {
  const paretoDiagramWidget = this.page.widget.paretoDiagram();

  this.completed.widget.setCommonFields(fields);

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => paretoDiagramWidget.clickSubmitParetoDiagram(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
