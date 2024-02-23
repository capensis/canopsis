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
          v-html="description"
          class="pre-wrap"
        />
        <img
          v-if="image"
          :src="image"
          class="my-2"
          alt=""
          @click="showImageViewerModal"
        >
      </v-layout>
    </v-expand-transition>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { eventFilterActionsTypesImages } from '@/assets';

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
      const imageName = `${this.$i18n.locale.toUpperCase()}_${this.type}`;

      return eventFilterActionsTypesImages[`./${imageName}.png`]
        ?? eventFilterActionsTypesImages[`./${imageName}.svg`]
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
