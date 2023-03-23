<template lang="pug">
  v-layout.chart-widget-filters(wrap)
    c-quick-date-interval-field.mr-4(
      v-if="showInterval",
      :interval="interval",
      :min="minIntervalDate",
      :quick-ranges="quickRanges",
      @input="$emit('update:interval', $event)"
    )
    c-sampling-field.chart-widget-filters__sampling.mr-4(
      v-if="showSampling",
      :value="sampling",
      @input="$emit('update:sampling', $event)"
    )
    v-layout(v-if="showFilter", row, align-end)
      filter-selector.chart-widget-filters__filter-selector.mr-4(
        :label="$t('settings.selectAFilter')",
        :filters="userFilters",
        :locked-filters="widgetFilters",
        :locked-value="lockedFilter",
        :value="filters",
        :disabled="filterDisabled",
        clearable,
        hide-details,
        @input="$emit('update:filters', $event)"
      )
      filters-list-btn(
        :widget-id="widgetId",
        :addable="filterAddable",
        :editable="filterEditable",
        private,
        with-alarm,
        with-entity,
        with-pbehavior
    )
</template>

<script>
import { METRICS_QUICK_RANGES } from '@/constants';

import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';
import FilterSelector from '@/components/other/filter/filter-selector.vue';

export default {
  components: { FilterSelector, FiltersListBtn },
  props: {
    widgetId: {
      type: String,
      required: true,
    },
    interval: {
      type: Object,
      required: false,
    },
    sampling: {
      type: String,
      required: false,
    },
    filters: {
      type: [String, Array],
      required: false,
    },
    userFilters: {
      type: Array,
      required: false,
    },
    widgetFilters: {
      type: Array,
      required: true,
    },
    lockedFilter: {
      type: [String, Array],
      required: false,
    },
    minIntervalDate: {
      type: Number,
      required: false,
    },
    showFilter: {
      type: Boolean,
      default: false,
    },
    showSampling: {
      type: Boolean,
      default: false,
    },
    showInterval: {
      type: Boolean,
      default: false,
    },
    filterAddable: {
      type: Boolean,
      default: false,
    },
    filterEditable: {
      type: Boolean,
      default: false,
    },
    filterDisabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    quickRanges() {
      return Object.values(METRICS_QUICK_RANGES);
    },
  },
};
</script>

<style lang="scss">
.chart-widget-filters {
  &__sampling {
    max-width: 200px;
  }

  &__filter-selector {
    max-width: 300px;
  }
}
</style>
