// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  clickUserDropDown() {
    return this.customClick('@userTopBarDropdownButton');
  },

  clickUserProfileButton() {
    const { userProfileButton } = this.elements;

    return this.clickUserDropDownItem(userProfileButton.locateStrategy, userProfileButton.selector);
  },

  clickLogoutButton() {
    return this.customClick('@logoutButton');
  },

  clickUserDropDownItem(locateStategy, selector) {
    this.api.element(locateStategy, selector, ({ status }) => {
      if (status === 0) {
        this.clickUserDropDown();
      }

      return this.customClick(selector);
    });

    return this;
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
