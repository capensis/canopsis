// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  clickSubmitAlarms() {
    return this.customClick('@submitAlarms');
  },

  clickAdvancedSettings() {
    return this.customClick('@advancedSettings')
      .defaultPause();
  },
};

module.exports = {
  elements: {
    advancedSettings: sel('advancedSettings'),
    submitAlarms: sel('submitAlarms'),
  },
  commands: [commands],
};
