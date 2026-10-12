<template>
  <v-layout>
    <test-suite-videos :videos="testSuite.videos" />
  </v-layout>
</template>

<script>
import TestSuiteVideos from './test-suite-videos.vue';

export default {
  components: { TestSuiteVideos },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
};
</script>
