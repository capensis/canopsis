<template>
  <c-layout class="kpi-widget-filters" wrap>
    <c-quick-date-interval-field
      v-if="showInterval"
      :interval="interval"
      :min="minIntervalDate"
      :quick-ranges="quickRanges"
      class="kpi-widget-filters__interval"
      short
      @update:interval="$emit('update:interval', $event)"
    />
    <c-sampling-field
      v-if="showSampling"
      :value="sampling"
      class="kpi-widget-filters__sampling"
      @update:value="$emit('update:sampling', $event)"
    />
    <c-layout
      v-if="showFilter"
      align="end"
    >
      <filter-selector
        :model-value="filters"
        :label="$t('settings.selectAFilter')"
        :filters="userFilters"
        :locked-filters="widgetFilters"
        :locked-value="lockedFilters"
        class="kpi-widget-filters__filter-selector mr-4"
        clearable
        hide-details
        @update:model-value="$emit('update:filters', $event)"
      />
      <filters-list-btn
        :widget-id="widgetId"
        addable
        editable
        private
        with-entity
      />
    </c-layout>
  </c-layout>
</template>

<script>
import { METRICS_QUICK_RANGES } from '@/constants';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

export default {
  components: {
    FilterSelector,
    FiltersListBtn,
  },

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
    lockedFilters: {
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
  },

  emits: ['update:interval', 'update:sampling', 'update:filters'],

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

  &__filter-selector,
  &__interval {
    max-width: 200px;
  }
}
</style>
