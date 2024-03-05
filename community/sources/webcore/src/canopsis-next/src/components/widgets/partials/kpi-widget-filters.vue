<template>
  <v-layout
    class="kpi-widget-filters"
    wrap
  >
    <c-quick-date-interval-field
      v-if="showInterval"
      :interval="interval"
      :min="minIntervalDate"
      :quick-ranges="quickRanges"
      class="kpi-widget-filters__interval"
      short
      @input="$emit('update:interval', $event)"
    />
    <c-sampling-field
      v-if="showSampling"
      :value="sampling"
      class="kpi-widget-filters__sampling"
      @input="$emit('update:sampling', $event)"
    />
    <v-layout
      v-if="showFilter"
      align-end
    >
      <filter-selector
        :label="$t('settings.selectAFilter')"
        :filters="userFilters"
        :locked-filters="widgetFilters"
        :locked-value="lockedFilter"
        :value="filters"
        :disabled="filterDisabled"
        class="kpi-widget-filters__filter-selector mr-4"
        clearable
        hide-details
        @input="$emit('update:filters', $event)"
      />
      <filters-list-btn
        v-if="filterAddable || filterEditable"
        :widget-id="widgetId"
        :addable="filterAddable"
        :editable="filterEditable"
        private
        with-entity
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { METRICS_QUICK_RANGES } from '@/constants';

import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';
import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';

export default {
  components: { FilterSelector, FiltersListBtn },
  props: {
    widgetId: {
      type: String,
      required: false,
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
      default: () => [],
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

<style lang="scss" scoped>
.kpi-widget-filters {
  column-gap: 24px;

  &__sampling {
    max-width: 100px;
  }

  &__filter-selector, &__interval {
    max-width: 200px;
  }
}
</style>
