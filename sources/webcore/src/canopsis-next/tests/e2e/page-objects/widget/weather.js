// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitWeather() {
    return this.customClick('@submitWeather');
  },

  clickTemplateWeatherItem() {
    return this.customClick('@templateWeatherItem');
  },

  clickTemplateModal() {
    return this.customClick('@templateModal');
  },

  clickTemplateEntities() {
    return this.customClick('@templateEntities');
  },

  clickWeatherItem(id) {
    return this.customClick(this.el('@templateEntities', id));
  },

  clickWeatherItemHelp(id) {
    return this.customClick(this.el('@weatherItemHelpButton', id));
  },

  clickWeatherItemPause(id) {
    return this.customClick(this.el('@weatherItemPauseButton', id));
  },

  clickWeatherItemSeeAlarms(id) {
    return this.customClick(this.el('@weatherItemPauseButton', id));
  },

  el,
};

module.exports = {
  elements: {
    submitWeather: sel('submitWeather'),

    weather: sel('weather'),
    weatherItem: sel('weatherItem-%s'),
    weatherItemHelpButton: `${sel('weatherItem-%s')} ${sel('weatherHelpButton')}`,
    weatherItemPauseButton: `${sel('weatherItem-%s')} ${sel('weatherPauseButton')}`,
    weatherItemSeeAlarmsButton: `${sel('weatherItem-%s')} ${sel('weatherSeeAlarmsButton')}`,

    templateWeatherItem: `${sel('widgetTemplateWeatherItem')} ${sel('showEditorModalButton')}`,
    templateModal: `${sel('widgetTemplateModal')} ${sel('showEditorModalButton')}`,
    templateEntities: `${sel('widgetTemplateEntities')} ${sel('showEditorModalButton')}`,
  },
  commands: [commands],
};
