<template>
  <v-container
    class="pa-0"
    fluid
    grid-list-sm
  >
    <v-layout wrap>
      <v-flex
        v-for="video in videoSources"
        :key="video"
        xs6
      >
        <video
          :src="video"
          class="d-block"
          width="100%"
          controls
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
