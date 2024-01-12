<template>
  <v-layout
    class="remediation-statistics-filters my-4"
    wrap
  >
    <c-quick-date-interval-field
      class="mr-4"
      v-field="query.interval"
      :min="minDate"
    />
    <c-sampling-field
      class="remediation-statistics-filters__sampling mr-4"
      v-field="query.sampling"
    />
    <remediation-statistics-show-mode-field
      class="remediation-statistics-filters__type mr-4"
      v-field="query.type"
    />
    <remediation-statistics-instructions-field
      class="remediation-statistics-filters__instruction"
      v-field="query.instruction"
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
