<template lang="pug">
  div.kpi-sli-toolbar
    v-layout.ml-4.my-4(wrap)
      c-quick-date-interval-field.mr-4(v-field="query.interval", :min="minDate")
      c-sampling-field.mr-4.kpi-sli-toolbar__sampling(v-field="query.sampling")
      kpi-sli-show-mode-field.mr-4.kpi-sli-toolbar__show-mode(v-field="query.type")
      c-filter-field.kpi-sli-toolbar__filters(v-field="query.filter")
</template>

<script>
import KpiSliShowModeField from './kpi-sli-show-mode-field.vue';

export default {
  components: { KpiSliShowModeField },
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
