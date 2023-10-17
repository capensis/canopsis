<template>
  <v-container
    class="pa-0"
    fluid
    grid-list-sm
  >
    <v-layout wrap>
      <v-flex
        xs4
        v-for="image in images"
        :key="image.src"
      >
        <v-img
          class="cursor-pointer"
          :src="image.src"
          aspect-ratio="1"
          @click="showImagesModal(image.src)"
        />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { MODALS } from '@/constants';

import { getTestSuiteFileUrl } from '@/helpers/entities/junit/url';

export default {
  props: {
    screenshots: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    images() {
      return this.screenshots.map(id => ({ src: getTestSuiteFileUrl(id) }));
    },
  },
  methods: {
    showImagesModal(screenshot) {
      this.$modals.show({
        name: MODALS.imagesViewer,
        config: {
          images: this.images,
          active: screenshot,
        },
      });
    },
  },
};
</script>
