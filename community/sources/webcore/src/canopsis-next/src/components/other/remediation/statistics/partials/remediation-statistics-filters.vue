<template lang="pug">
  div.kpi-sli-toolbar
    v-layout.my-4(row)
      c-quick-date-interval-field.mr-4(v-field="query.interval", :min="minDate")
      c-sampling-field.mr-4(v-field="query.sampling")
      remediation-statistics-show-mode-field.mr-4(v-field="query.type")
      remediation-statistics-instructions-field(v-field="query.instruction")
</template>

<script>
import RemediationStatisticsShowModeField from './remediation-statistics-show-mode-field.vue';
import RemediationStatisticsInstructionsField from './remediation-statistics-instructions-field.vue';

export default {
  components: {
    RemediationStatisticsShowModeField,
    RemediationStatisticsInstructionsField,
  },
  model: {
    prop: 'query',
    event: 'input',
  },
  props: {
    query: {
      type: Object,
      required: true,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },
};
</script>

<style scoped lang="scss">
.kpi-sli-toolbar {
  &__sampling {
    max-width: 200px;
  }

  &__show-mode {
    max-width: 150px;
  }

  &__filters {
    max-width: 200px;
  }
}
</style>
