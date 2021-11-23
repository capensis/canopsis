<template lang="pug">
  div.kpi-alarms-toolbar
    v-layout.ml-4.my-4(wrap)
      c-quick-date-interval-field.mr-4(v-field="query.interval")
      c-sampling-field.mr-4.kpi-alarms-toolbar__sampling(v-field="query.sampling")
      c-filters-field.mr-4.kpi-alarms-toolbar__filters(v-field="query.filter")
      c-alarm-metric-parameters-field.kpi-alarms-toolbar__parameters(v-field="query.parameters")
</template>

<script>
export default {
  model: {
    prop: 'query',
    event: 'input',
  },
  props: {
    query: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style scoped lang="scss">
.kpi-alarms-toolbar {
  &__sampling {
    max-width: 200px;
  }

  &__filters {
    max-width: 200px;
  }

  &__parameters {
    max-width: 300px;
  }
}
</style>
