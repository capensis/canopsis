<template>
  <c-layout>
    <test-suite-screenshots :screenshots="testSuite.screenshots" />
  </c-layout>
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
