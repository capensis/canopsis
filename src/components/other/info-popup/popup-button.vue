<template lang="pug">
  span.info-popup-button(v-if="popupData", @click.stop="showPopup")
    v-icon info
</template>

<script>
import Handlebars from 'handlebars';
import getProp from 'lodash/get';

import popupComponentMixin from '@/mixins/popup';

/**
 * Button to display info popup
 *
 * @prop {String} [columnName] - Name of the column
 * @prop {Object} [alarm] - Object representing the alarm
 */
export default {
  mixins: [
    popupComponentMixin,
  ],
  props: {
    columnName: {
      type: String,
      required: true,
    },
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    popupData() {
      const popups = getProp(this.widget, 'widget.popup', []);

      return popups.find(popup => popup.column === this.columnName);
    },
    textContent() {
      const template = Handlebars.compile(this.popupData.template);
      const context = { alarm: this.alarm.v };

      return template(context);
    },
  },
  methods: {
    showPopup() {
      this.addInfoPopup({
        text: this.textContent,
        autoClose: false,
      });
    },
  },
};
</script>

<style scoped>
  .info-popup-button {
    cursor: pointer;
  }
</style>
