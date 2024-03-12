<template>
  <v-layout
    class="remediation-statistics-filters my-4"
    wrap
  >
    <c-quick-date-interval-field
      v-field="query.interval"
      :min="minDate"
      class="mr-4"
    />
    <c-sampling-field
      v-field="query.sampling"
      class="remediation-statistics-filters__sampling mr-4"
    />
    <remediation-statistics-show-mode-field
      v-field="query.type"
      class="remediation-statistics-filters__type mr-4"
    />
    <remediation-statistics-instructions-field
      v-field="query.instruction"
      class="remediation-statistics-filters__instruction"
    />
  </v-layout>
</template>

<script>
import RemediationStatisticsShowModeField from '../form/fields/remediation-statistics-show-mode-field.vue';
import RemediationStatisticsInstructionsField from '../form/fields/remediation-statistics-instructions-field.vue';

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

<style lang="scss">
.remediation-statistics-filters {
  &__sampling {
    max-width: 200px;
  }

  &__type {
    max-width: 150px;
  }

  &__instruction {
    max-width: 400px;
  }
}
</style>
