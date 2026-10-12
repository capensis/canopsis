<template>
  <div>
    <test-suite-history-panel
      v-for="historyTestSuite in testSuiteHistory"
      :key="historyTestSuite._id"
      :test-suite="historyTestSuite"
      class="mb-2"
    />
  </div>
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
      default: () => [],
    },
  },
};
</script>
