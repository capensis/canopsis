<template>
  <c-layout column>
    <c-information-block-row
      :label="$t('testSuite.xmlFeed')"
      :value="testSuite.xml_feed"
    />
    <c-information-block-row
      :label="$t('common.name')"
      :value="testSuite.name"
    />
    <c-information-block-row
      :label="$t('testSuite.hostname')"
      :value="testSuite.hostname"
    />
    <c-information-block-row :label="$t('testSuite.lastUpdate')">
      {{ convertDateToString(testSuite.last_update, 'testSuiteFormat') }}
    </c-information-block-row>
    <c-information-block-row :label="$t('common.timeTaken')">
      <span v-if="testSuite.time">
        {{ convertNumberToFixedString(testSuite.time) }}{{ $constants.TIME_UNITS.second }}
      </span>
    </c-information-block-row>
    <v-row class="mt-4">
      <c-layout column>
        <c-information-block-row
          :label="$t('testSuite.totalTests')"
          :value="testSuite.total"
        />
        <test-suite-summary-status-row
          :label="$t('testSuite.disabledTests')"
          :total="testSuite.total"
          :count="testSuite.disabled"
        />
        <test-suite-summary-status-row
          :label="$tc('common.error', 2)"
          :total="testSuite.total"
          :count="testSuite.errors"
        />
        <test-suite-summary-status-row
          :label="$t('common.failures')"
          :total="testSuite.total"
          :count="testSuite.failures"
        />
        <test-suite-summary-status-row
          :label="$t('common.skipped')"
          :total="testSuite.total"
          :count="testSuite.skipped"
        />
      </c-layout>
      <v-col cols="4">
        <test-suite-status-pie-chart :statuses="testSuiteStatuses" />
      </v-col>
    </v-row>
  </c-layout>
</template>

<script>
import { convertDateToString } from '@/helpers/date/date';
import { convertNumberToFixedString } from '@/helpers/string';

import TestSuiteStatusPieChart from './test-suite-status-pie-chart.vue';
import TestSuiteSummaryStatusRow from './test-suite-summary-status-row.vue';

export default {
  components: {
    TestSuiteStatusPieChart,
    TestSuiteSummaryStatusRow,
  },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
  computed: {
    passed() {
      const {
        disabled, total, failures, skipped, errors,
      } = this.testSuite;

      return total - (disabled + failures + skipped + errors);
    },

    testSuiteStatuses() {
      return {
        skipped: this.testSuite.skipped,
        failed: this.testSuite.failures,
        error: this.testSuite.errors,
        passed: this.passed,
      };
    },
  },
  methods: {
    convertDateToString,
    convertNumberToFixedString,
  },
};
</script>
