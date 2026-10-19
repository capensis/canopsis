// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsCalendar() {
    return this.customClick('@submitStatsCalendar');
  },

  clickCriticityLevels() {
    return this.customClick('@criticityLevels');
  },

  clickCriticityLevelsMinor() {
    return this.customClick('@criticityLevelsMinor');
  },

  clearCriticityLevelsMinor() {
    return this.customClearValue('@criticityLevelsMinor');
  },

  setCriticityLevelsMinor(value = 20) {
    return this.customSetValue('@criticityLevelsMinor', value);
  },

  clickCriticityLevelsMajor() {
    return this.customClick('@criticityLevelsMajor');
  },

  clearCriticityLevelsMajor() {
    return this.customClearValue('@criticityLevelsMajor');
  },

  setCriticityLevelsMajor(value = 30) {
    return this.customSetValue('@criticityLevelsMajor', value);
  },

  clickCriticityLevelsCritical() {
    return this.customClick('@criticityLevelsCritical');
  },

  clearCriticityLevelsCritical() {
    return this.customClearValue('@criticityLevelsCritical');
  },

  setCriticityLevelsCritical(value = 40) {
    return this.customSetValue('@criticityLevelsCritical', value);
  },

  setConsiderPbehaviors(checked = false) {
    return this.getAttribute('@widgetConsiderPbehaviorsInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@widgetConsiderPbehaviors');
      }
    });
  },

  clickColorSelector() {
    return this.customClick('@levelsColorsSelector');
  },

  clickColorPickerButton(level) {
    return this.customClick(this.el('@levelColorPicker', level));
  },

  el,
};

module.exports = {
  elements: {
    submitStatsCalendar: sel('submitStatsCalendarButton'),

    criticityLevels: sel('widgetCriticityLevels'),
    criticityLevelsMinor: sel('criticityLevelsMinor'),
    criticityLevelsMajor: sel('criticityLevelsMajor'),
    criticityLevelsCritical: sel('criticityLevelsCritical'),

    widgetConsiderPbehaviorsInput: `${sel('widgetConsiderPbehaviors')} input${sel('switcherField')}`,
    widgetConsiderPbehaviors: `div${sel('widgetConsiderPbehaviors')}`,

    levelsColorsSelector: `div${sel('levelsColorsSelector')}`,
    levelColorPicker: `${sel('levelsColor-%s')} ${sel('showColorPickerButton')}`,
  },
  commands: [commands],
};
