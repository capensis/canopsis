// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  clickSubmitText() {
    return this.customClick('@submitText');
  },

  clickCreateTemplate() {
    return this.customClick('@createTemplate');
  },

  clickEditTemplate() {
    return this.customClick('@editTemplate');
  },

  clickDeleteTemplate() {
    return this.customClick('@deleteTemplate');
  },

  clickStats() {
    return this.customClick('@textStats');
  },
};

module.exports = {
  elements: {
    submitText: sel('submitText'),

    createTemplate: `${sel('widgetTestTemplate')} ${sel('createButton')}`,
    editTemplate: `${sel('widgetTestTemplate')} ${sel('editButton')}`,
    deleteTemplate: `${sel('widgetTestTemplate')} ${sel('deleteButton')}`,

    textStats: sel('textWidgetStats'),
  },
  commands: [commands],
};
