<template lang="pug">
  div
    v-container(v-for="(popup, index) in popups" :key="`info-popup-${index}`")
      v-card
        v-container
          v-layout(class="text-md-center")
            v-flex
              v-btn(color="error", @click="deleteItem(index)") Delete
          v-layout
            v-flex(:justify-center="true")
              v-text-field(placeholder="Column", v-model="popup.column")
          v-layout
            v-flex(:justify-center="true")
              v-text-field(placeholder="Template", :multi-line="true", v-model="popup.template")
    v-btn(color="success", @click="add") Add
    v-btn(color="info", @click="save") Save

</template>

<script>
import alarmInfoPopupMixin from '@/mixins/widget';

export default {
  mixins: [
    alarmInfoPopupMixin,
  ],
  data() {
    return {
      popups: [],
    };
  },
  computed: {
    widget() {
      return this.getWidget({ widgetXType: 'listalarm' });
    },
  },
  watch: {
    widget() {
      if (this.widget) {
        this.popups = this.widget.widget.popup;
      }
    },
  },
  methods: {
    deleteItem(index) {
      this.popups.splice(index, 1);
    },
    add() {
      this.popups.push({
        column: '',
        template: '',
      });
    },
    async save() {
      await this.saveWidget({ widgetWrapper: this.widget });
    },
  },
};
</script>
