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

module.exports.generateTemporaryRight = (prefix = '', suffix = uid()) => ({
  ID: `${prefix}ID-${suffix}`,
  description: `${prefix}description-${suffix}`,
});

module.exports.generateTemporaryRole = (prefix = '', suffix = uid()) => ({
  name: `${prefix}name-${suffix}`,
  description: `${prefix}description-${suffix}`,
});
