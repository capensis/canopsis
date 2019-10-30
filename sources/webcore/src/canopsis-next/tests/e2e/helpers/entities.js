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
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryAlarmsWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Alarms-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryContextWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Context-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryStatsTableWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}StatsTable-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryTextWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Text-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryStatsCalendarWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Stats-Calendar-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryStatsHistogramWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Stats-Histogram-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryStatsCurvesWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Stats-Curves-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryStatsNumberWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Stats-Number-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});

module.exports.generateTemporaryParetoDiagramWidget = (prefix = '', suffix = uid()) => ({
  title: `${prefix}Pareto-Diagram-widget-${suffix}`,
  row: `${prefix}row-${suffix}`,
  size: {
    sm: 12,
    md: 12,
    lg: 12,
  },
});
