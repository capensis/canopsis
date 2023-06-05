<template lang="pug">
  div
    v-layout(row, align-center)
      span {{ message }}
      v-btn(icon, small, @click="toggleDescriptionOpened")
        v-icon help
    v-expand-transition(v-if="opened")
      v-layout(column)
        div.pre-wrap {{ description }}
        img(:src="image", @click="showImageViewerModal")
</template>

<script>
import { MODALS } from '@/constants';

export default {
  props: {
    type: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    message() {
      return this.$t(`eventFilter.actionsTypes.${this.type}.message`);
    },

    description() {
      return this.$t(`eventFilter.actionsTypes.${this.type}.description`);
    },

    image() {
      // eslint-disable-next-line import/no-dynamic-require,global-require
      return require(`@/assets/event-filter-actions-types/${this.$i18n.locale.toUpperCase()}_${this.type}.png`);
    },
  },
  methods: {
    toggleDescriptionOpened() {
      this.opened = !this.opened;
    },

    showImageViewerModal() {
      this.$modals.show({
        name: MODALS.imageViewer,
        config: {
          src: this.image,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
img {
  width: 100%;
  cursor: pointer;
}
</style>
