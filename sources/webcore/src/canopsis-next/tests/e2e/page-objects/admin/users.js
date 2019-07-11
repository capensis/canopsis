// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/users`;
  },
  elements: {
    dataTable: '.v-datatable',
    addButton: sel('addButton'),
  },
  commands: [commands],
};
