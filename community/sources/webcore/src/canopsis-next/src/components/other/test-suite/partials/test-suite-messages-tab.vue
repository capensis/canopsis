<template>
  <v-layout column>
    <template v-if="hasMessages">
      <test-suite-message-panel
        class="mb-2"
        v-if="testSuite.system_err"
        :value="testSuite.system_err"
        :title="$t('testSuite.systemError')"
        :label="$t('testSuite.systemErrorMessage')"
        :file-name="systemErrorFileName"
        :color="$config.COLORS.testSuiteStatuses.failed"
      />
      <test-suite-message-panel
        v-if="testSuite.system_out"
        :value="testSuite.system_out"
        :title="$t('testSuite.systemOut')"
        :label="$t('testSuite.systemOutMessage')"
        :file-name="systemOutFileName"
        :color="$config.COLORS.testSuiteStatuses.error"
      />
    </template>
    <template v-else>
      <div>{{ testSuite.artifact_match_err }}</div>
      <div>{{ $t('testSuite.noData') }}</div>
    </template>
  </v-layout>
</template>

<script>
import TestSuiteMessagePanel from './test-suite-message-panel.vue';

export default {
  components: {
    TestSuiteMessagePanel,
  },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
  computed: {
    systemErrorFileName() {
      return `${this.testSuite.name}_system_error`;
    },

    systemOutFileName() {
      return `${this.testSuite.name}_system_out`;
    },

    hasMessages() {
      return this.testSuite.system_err || this.testSuite.system_out;
    },
  },
};
</script>
