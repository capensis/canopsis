<template lang="pug">
  v-container.pa-0(fluid, grid-list-sm)
    v-layout(row, wrap)
      v-flex(v-for="video in videoSources", :key="video", xs6)
        video.d-block(controls, :src="video", width="100%")
          source(:src="video")
</template>

<script>
import { getTestSuiteFileUrl } from '@/helpers/url';

export default {
  props: {
    videos: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    videoSources() {
      return this.videos.map(getTestSuiteFileUrl);
    },
  },
};
</script>
