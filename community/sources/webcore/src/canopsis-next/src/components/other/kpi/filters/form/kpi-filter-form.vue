<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    kpi-filter-patterns-form.mt-2(v-field="form.patterns")
</template>

<script>
import KpiFilterPatternsForm from './partials/kpi-filter-patterns-form.vue';

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
