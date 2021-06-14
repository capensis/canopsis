<template lang="pug">
  v-container.pa-0(fluid, grid-list-sm)
    v-layout(row, wrap)
      v-flex(xs4, v-for="image in images", :key="image.src")
        v-img.cursor-pointer(
          :src="image.src",
          aspect-ratio="1",
          @click="showImagesModal(image.src)"
        )
</template>

<script>
import { MODALS } from '@/constants';

import { getTestSuiteFileUrl } from '@/helpers/url';

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
