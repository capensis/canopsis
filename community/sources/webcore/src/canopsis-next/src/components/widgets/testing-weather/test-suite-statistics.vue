<template>
  <v-row>
    <v-col cols="4">
      {{ $t('common.tests') }}
    </v-col>
    <v-col cols="8">
      <test-suite-statistic-row
        :label="$t('common.total')"
        :count="testSuite.total"
      />
      <test-suite-statistic-row
        :label="$tc('common.error', 2)"
        :count="testSuite.errors"
      />
      <test-suite-statistic-row
        :label="$t('common.failures')"
        :count="testSuite.failures"
      />
      <test-suite-statistic-row
        :label="$t('common.skipped')"
        :count="testSuite.skipped"
      />
    </v-col>
  </v-row>
</template>

<script>
import TestSuiteStatisticRow from './test-suite-statistic-row.vue';

export default {
  components: {
    TestSuiteStatisticRow,
  },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
};
</script>
