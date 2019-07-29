const uid = require('uid');

module.exports.generateTemporaryView = (prefix = '', suffix = uid()) => ({
  name: `${prefix}name-${suffix}`,
  title: `${prefix}title-${suffix}`,
  description: `${prefix}description-${suffix}`,
  enabled: true,
  tags: [`${prefix}tag-${suffix}`],
  group: `${prefix}group-${suffix}`,
});
