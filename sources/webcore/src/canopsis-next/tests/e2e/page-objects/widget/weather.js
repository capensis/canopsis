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

    templateWeatherItem: `${sel('widgetTemplateWeatherItem')} ${sel('showEditorModalButton')}`,
    templateModal: `${sel('widgetTemplateModal')} ${sel('showEditorModalButton')}`,
    templateEntities: `${sel('widgetTemplateEntities')} ${sel('showEditorModalButton')}`,
  },
  commands: [commands],
};
