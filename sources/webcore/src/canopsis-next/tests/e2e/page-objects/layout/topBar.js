// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  clickUserDropdown() {
    return this.customClick('@userTopBarDropdownButton');
  },

  clickUserProfileButton() {
    return this.customClick('@userProfileButton');
  },

  clickLogoutButton() {
    return this.customClick('@logoutButton');
  },
};

module.exports = {
  elements: {
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    userProfileButton: sel('userProfileButton'),
    logoutButton: sel('logoutButton'),
  },
  commands: [commands],
};
