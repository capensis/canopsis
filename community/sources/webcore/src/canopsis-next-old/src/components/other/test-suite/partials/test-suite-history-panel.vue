<template>
  <c-collapse-panel
    :title="title"
    :color="color"
  >
    <test-suite-tabs :test-suite-id="testSuite._id" />
  </c-collapse-panel>
</template>

<script>
import { TEST_SUITE_COLORS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import TestSuiteTabs from './test-suite-tabs.vue';

export default {
  components: { TestSuiteTabs },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
  computed: {
    title() {
      return convertDateToString(this.testSuite.last_update || this.testSuite.created, 'testSuiteFormat');
    },

    color() {
      return TEST_SUITE_COLORS[this.testSuite.state];
    },
  },
};
</script>
