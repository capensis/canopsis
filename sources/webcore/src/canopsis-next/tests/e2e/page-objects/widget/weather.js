// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  clickSubmitWeather() {
    return this.customClick('@submitWeather');
  },

  clickAdvancedSettings() {
    return this.customClick('@advancedSettings')
      .defaultPause();
  },
};

module.exports = {
  elements: {
    advancedSettings: sel('advancedSettings'),
    submitWeather: sel('submitWeather'),
  },
  commands: [commands],
};
