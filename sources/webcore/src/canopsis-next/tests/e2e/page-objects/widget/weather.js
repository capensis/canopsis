// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  clickSubmitWeather() {
    return this.customClick('@submitWeather');
  },

  clickTemplateWeatherItem() {
    return this.customClick('@templateWeatherItem')
      .defaultPause();
  },

  clickTemplateModal() {
    return this.customClick('@templateModal')
      .defaultPause();
  },

  clickTemplateEntities() {
    return this.customClick('@templateEntities')
      .defaultPause();
  },
};

module.exports = {
  elements: {
    submitWeather: sel('submitWeather'),

    templateWeatherItem: `${sel('widgetTemplateWeatherItem')} ${sel('showEditButton')}`,
    templateModal: `${sel('widgetTemplateModal')} ${sel('showEditButton')}`,
    templateEntities: `${sel('widgetTemplateEntities')} ${sel('showEditButton')}`,
  },
  commands: [commands],
};
