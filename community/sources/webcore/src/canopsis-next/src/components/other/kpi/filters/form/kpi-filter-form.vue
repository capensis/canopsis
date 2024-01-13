<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <kpi-filter-patterns-form
      class="mt-2"
      v-field="form.patterns"
    />
  </v-layout>
</template>

<script>
import KpiFilterPatternsForm from './kpi-filter-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: { KpiFilterPatternsForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
