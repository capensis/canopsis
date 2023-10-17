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
          class="d-block"
          controls
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
