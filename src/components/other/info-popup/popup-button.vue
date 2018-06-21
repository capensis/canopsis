<template lang="pug">
  span.info-popup-button(v-if="popupData", @click.stop="showPopup")
    v-icon info
</template>

<script>
import Handlebars from 'handlebars';
import getProp from 'lodash/get';

import widgetMixin from '@/mixins/widget';
import popupComponentMixin from '@/mixins/popup';

export default {
  mixins: [
    widgetMixin,
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
  },
  computed: {
    popupData() {
      const popups = getProp(this.getWidget({ widgetXType: 'listalarm' }), 'widget.popup', []);

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
