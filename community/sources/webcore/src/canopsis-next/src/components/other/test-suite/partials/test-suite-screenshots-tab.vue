<template lang="pug">
  v-layout.pa-3
    test-suite-screenshots(:screenshots="testSuite.screenshots")
</template>

<script>
import TestSuiteScreenshots from './test-suite-screenshots.vue';

export default {
  components: { TestSuiteScreenshots },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
};
</script>
