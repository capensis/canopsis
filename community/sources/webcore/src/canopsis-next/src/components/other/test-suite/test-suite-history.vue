<template lang="pug">
  div
    test-suite-history-panel.mb-2(
      v-for="testSuite in testSuiteHistory",
      :key="testSuite._id",
      :test-suite="testSuite"
    )
</template>

<script>
import TestSuiteHistoryPanel from './partials/test-suite-history-panel.vue';

export default {
  components: {
    TestSuiteHistoryPanel,
  },
  props: {
    testSuiteHistory: {
      type: Array,
      required: () => [],
    },
  },
};
</script>
