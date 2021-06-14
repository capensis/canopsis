<template lang="pug">
  test-suite-summary-row(:label="label")
    v-layout
      v-flex(xs2) {{ count }}
      v-flex(v-if="testsStatusPercent") {{ testsStatusPercent | percentage(precision) }}
</template>

<script>
import TestSuiteSummaryRow from './test-suite-summary-row.vue';

export default {
  components: { TestSuiteSummaryRow },
  props: {
    label: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      required: false,
    },
    total: {
      type: Number,
      required: false,
    },
    precision: {
      type: Number,
      default: 2,
    },
  },
  computed: {
    testsStatusPercent() {
      if (this.count === 0) {
        return this.count;
      }

      return this.count / this.total;
    },
  },
};
</script>
