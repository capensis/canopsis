const uid = require('uid');

module.exports.generateTemporaryView = (prefix = '', suffix = uid()) => ({
  name: `${prefix}name-${suffix}`,
  title: `${prefix}title-${suffix}`,
  description: `${prefix}description-${suffix}`,
  enabled: true,
  tags: [`${prefix}tag-${suffix}`],
  group: `${prefix}group-${suffix}`,
});

module.exports.generateTemporaryUser = (prefix = '', suffix = uid()) => ({
  username: `${prefix}username-${suffix}`,
  firstname: `${prefix}firstname-${suffix}`,
  lastname: `${prefix}lastname-${suffix}`,
  email: `${prefix}email-${suffix}@example.com`,
  password: `${prefix}password-${suffix}`,
});

module.exports.generateTemporaryRole = (prefix = '', suffix = uid()) => ({
  name: `${prefix}name-${suffix}`,
  description: `${prefix}description-${suffix}`,
});

module.exports.generateTemporaryWeatherWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Weather-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
});

module.exports.generateTemporaryAlarmsWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Alarms-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
});

module.exports.generateTemporaryContextWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Context-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
});

module.exports.generateTemporaryStatsTableWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}StatsTable-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
});

module.exports.generateTemporaryTextWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Text-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
});

module.exports.generateTemporaryStatsCalendarWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Stats-Calendar-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
});
