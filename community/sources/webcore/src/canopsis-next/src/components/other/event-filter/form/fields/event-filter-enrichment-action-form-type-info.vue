<template>
  <div>
    <v-layout align-center>
      <span>{{ message }}</span>
      <v-btn
        icon
        small
        @click="toggleDescriptionOpened"
      >
        <v-icon>help</v-icon>
      </v-btn>
    </v-layout>
    <v-expand-transition v-if="opened">
      <v-layout column>
        <div class="pre-wrap">
          {{ description }}
        </div>
        <img
          :src="image"
          alt=""
          @click="showImageViewerModal"
        >
      </v-layout>
    </v-expand-transition>
  </div>
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
