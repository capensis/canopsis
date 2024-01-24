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
    <v-expand-transition>
      <v-layout
        v-show="opened"
        column
      >
        <div
          class="pre-wrap"
          v-html="description"
        />
        <img
          class="my-2"
          v-if="image"
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

const images = {};
const requireModule = require.context('../../../../../assets/event-filter-actions-types/');

requireModule.keys().forEach(key => images[key] = requireModule(key));

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
      return images[`./${this.$i18n.locale.toUpperCase()}_${this.type}.png`]
        ?? images[`./${this.$i18n.locale.toUpperCase()}_${this.type}.svg`]
        ?? '';
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
