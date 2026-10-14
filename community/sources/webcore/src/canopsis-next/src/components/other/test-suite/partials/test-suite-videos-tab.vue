<template>
  <c-layout>
    <test-suite-videos :videos="testSuite.videos" />
  </c-layout>
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
