<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    kpi-filter-patterns-form(v-field="form.patterns")
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
