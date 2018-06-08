<template lang="pug">
  span(v-if="hasPopupForColumn", @click.stop="showPopup", class="info-popup-button")
    v-icon info
</template>

<script>
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
  },
  methods: {
    showPopup() {
      this.addInfoPopup({
        text: this.popupData.template,
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
