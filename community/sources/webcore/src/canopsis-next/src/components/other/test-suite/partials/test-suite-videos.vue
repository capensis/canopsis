<template>
  <v-container
    class="pa-0"
    fluid="fluid"
    grid-list-sm="grid-list-sm"
  >
    <v-layout wrap="wrap">
      <v-flex
        v-for="video in videoSources"
        :key="video"
        xs6="xs6"
      >
        <video
          class="d-block"
          controls="controls"
          :src="video"
          width="100%"
        >
          <source :src="video">
        </video>
      </v-flex>
    </v-layout>
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
