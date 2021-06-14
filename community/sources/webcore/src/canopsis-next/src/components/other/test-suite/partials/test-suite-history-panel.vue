<template lang="pug">
  div.test-suite-history
    v-expansion-panel
      v-expansion-panel-content(:style="{ backgroundColor: color }", lazy)
        v-icon(slot="actions", color="white") $vuetify.icons.expand
        span.white--text(slot="header") {{ title }}
        test-suite-tabs(:test-suite-id="testSuite._id")
</template>

<script>
import { TEST_SUITE_COLORS } from '@/constants';

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
      const timestamp = this.testSuite.last_update || this.testSuite.created;

      return this.$options.filters.date(timestamp, 'testSuiteFormat', true);
    },

    color() {
      return TEST_SUITE_COLORS[this.testSuite.state];
    },
  },
};
</script>

<style lang="scss" scoped>
.test-suite-history {
  & /deep/ .v-expansion-panel {
    border-radius: 5px;
    overflow: hidden;
  }
}
</style>
