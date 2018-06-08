<template lang="pug">
  span(v-if="hasPopupForColumn", @click.stop="showPopup", class="info-popup-button")
    v-icon info
</template>

<script>
import widgetMixin from '@/mixins/widget';
import popupComponentMixin from '@/mixins/popup';
import Handlebars from 'handlebars';

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
    hasPopupForColumn() {
      return Boolean(this.popupData);
    },
    popupData() {
      const popups = this.getWidget({ widgetXType: 'listalarm' }).widget.popup;
      let data = null;

      popups.forEach((popup) => {
        if (popup.column === this.columnName) {
          data = popup;
        }
      });

      return data;
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
