<template>
  <v-container class="pa-0" fluid>
    <v-row density="compact">
      <v-col
        v-for="video in videoSources"
        :key="video"
        cols="6"
      >
        <video
          :src="video"
          class="d-block"
          width="100%"
          controls
        >
          <source :src="video">
        </video>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { getTestSuiteFileUrl } from '@/helpers/entities/junit/url';

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
