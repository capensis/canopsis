<template lang="pug">
  v-layout
    v-flex(xs12)
      test-suite-details(:test-suite="testSuite")
</template>

<script>
import TestSuiteDetails from './test-suite-details.vue';

export default {
  components: { TestSuiteDetails },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
};
</script>
