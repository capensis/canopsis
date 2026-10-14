<template>
  <v-container class="pa-0" fluid>
    <v-row density="compact">
      <v-col
        v-for="image in images"
        :key="image.src"
        cols="4"
      >
        <v-img
          :src="image.src"
          class="cursor-pointer"
          aspect-ratio="1"
          @click="showImagesModal(image.src)"
        />
      </v-col>
    </v-row>
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
