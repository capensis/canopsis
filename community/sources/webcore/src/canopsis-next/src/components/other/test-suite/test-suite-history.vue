<template>
  <div>
    <test-suite-history-panel
      class="mb-2"
      v-for="testSuite in testSuiteHistory"
      :key="testSuite._id"
      :test-suite="testSuite"
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
      required: () => [],
    },
  },
};
</script>
