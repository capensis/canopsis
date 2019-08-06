// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  el,

  clickRowGridSize() {
    return this.customClick('@rowGridSize');
  },
  setRow(value) {
    return this.customSetValue('@rowGridSizeCombobox', value)
      .customKeyup('@rowGridSizeCombobox', 'ENTER');
  },
  setSlider(slider, value) {
    return this.customSetValue(this.el('@slider', slider), value);
  },
};

module.exports = {
  elements: {
    rowGridSize: sel('rowGridSize'),
    rowGridSizeCombobox: sel('rowGridSizeCombobox'),
    slider: sel('slider-%s'),
  },
  commands: [commands],
};
