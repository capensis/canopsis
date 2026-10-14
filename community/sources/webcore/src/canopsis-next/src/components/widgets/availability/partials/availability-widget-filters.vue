<template>
  <c-layout class="availability-widget-filters col-gap-6 row-gap-3" wrap>
    <c-availability-advanced-search
      :searches="searches"
      class="mt-0 availability-widget-filters__search"
      @submit="$emit('submit:search', $event)"
      @reset="$emit('reset:search')"
      @toggle-pin:search="$emit('toggle-pin:search', $event)"
      @remove:search="$emit('remove:search', $event)"
    />
    <c-quick-date-interval-field
      v-if="showInterval"
      :interval="interval"
      :quick-ranges="quickRanges"
      class="availability-widget-filters__interval"
      short
      with-hours
      allow-future
      @update:interval="$emit('update:interval', $event)"
    />
    <c-layout
      v-if="showFilter"
      class="availability-widget-filters__filter-wrapper"
      align="end"
    >
      <filter-selector
        :label="$t('settings.selectAFilter')"
        :filters="userFilters"
        :locked-filters="widgetFilters"
        :locked-value="lockedFilter"
        :model-value="filters"
        :disabled="filterDisabled"
        class="availability-widget-filters__filter-selector mr-4"
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

    <availability-display-parameter-field
      :model-value="displayParameter"
      :label="$t('common.value')"
      class="availability-widget-filters__parameter"
      @update:model-value="$emit('update:display-parameter', $event)"
    />
    <availability-show-type-field
      :model-value="type"
      :label="$t('common.show')"
      class="availability-widget-filters__show-type"
      @update:model-value="$emit('update:type', $event)"
    />
    <c-enabled-field
      :model-value="trend"
      :label="$t('settings.showTrend')"
      class="mt-0"
      color="primary"
      hide-details
      inline
      no-margin
      @update:model-value="$emit('update:trend', $event)"
    />
    <availability-value-filter-field
      v-model="localValueFilter"
      :show-type="type"
      :max-seconds="maxValueFilterSeconds"
      class="availability-widget-filters__value-filter"
      @update:model-value="handleUpdateValueFilter"
    />
    <c-action-btn
      v-if="showExport"
      :loading="exporting"
      :tooltip="$t('settings.exportAsCsv')"
      icon="cloud_download"
      @click="$emit('export')"
    />
  </c-layout>
</template>

<script>
import { debounce } from 'lodash';
import { computed, ref, watch } from 'vue';

import { AVAILABILITY_QUICK_RANGES } from '@/constants';

import AvailabilityDisplayParameterField from '@/components/other/availability/form/fields/availability-display-parameter-field.vue';
import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';
import AvailabilityValueFilterField from '@/components/other/availability/form/fields/availability-value-filter-field.vue';
import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

export default {
  components: {
    AvailabilityValueFilterField,
    AvailabilityShowTypeField,
    AvailabilityDisplayParameterField,
    FilterSelector,
    FiltersListBtn,
  },

  props: {
    widgetId: {
      type: String,
      required: false,
    },
    query: {
      type: Object,
      required: false,
      default: () => ({}),
    },
    searches: {
      type: Array,
      default: () => [],
    },
    userFilters: {
      type: Array,
      required: false,
    },
    widgetFilters: {
      type: Array,
      default: () => [],
    },
    showInterval: {
      type: Boolean,
      default: false,
    },
    showFilter: {
      type: Boolean,
      default: false,
    },
    filterDisabled: {
      type: Boolean,
      default: false,
    },
    exporting: {
      type: Boolean,
      default: false,
    },
    maxValueFilterSeconds: {
      type: Number,
      required: false,
    },
    showExport: {
      type: Boolean,
      required: false,
    },
  },

  emits: [
    'submit:search',
    'reset:search',
    'toggle-pin:search',
    'remove:search',
    'update:interval',
    'update:filters',
    'update:display-parameter',
    'update:type',
    'update:trend',
    'update:value-filter',
    'export',
  ],

  setup(props, { emit }) {
    const localValueFilter = ref();
    const quickRanges = Object.values(AVAILABILITY_QUICK_RANGES);

    const interval = computed(() => props.query?.interval);
    const filters = computed(() => props.query?.filter);
    const lockedFilter = computed(() => props.query?.lockedFilter);
    const displayParameter = computed(() => props.query?.displayParameter);
    const type = computed(() => props.query?.showType);
    const trend = computed(() => props.query?.showTrend);
    const valueFilter = computed(() => props.query?.valueFilter);

    watch(
      valueFilter,
      () => {
        localValueFilter.value = valueFilter.value && { ...valueFilter.value };
      },
      { immediate: true },
    );

    const emitUpdateValueFilter = newValueFilter => emit('update:value-filter', newValueFilter);
    const debouncedEmitUpdateValueFilter = debounce(emitUpdateValueFilter, 1000);

    const handleUpdateValueFilter = (val) => {
      if (!val || val.value === valueFilter.value?.value) {
        emitUpdateValueFilter(val);
      } else {
        debouncedEmitUpdateValueFilter(val);
      }
    };

    return {
      localValueFilter,
      quickRanges,
      interval,
      filters,
      lockedFilter,
      displayParameter,
      type,
      trend,
      handleUpdateValueFilter,
    };
  },
};
</script>

<style lang="scss" scoped>
.availability-widget-filters {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;

  & > * {
    flex-grow: 0;
  }

  &__search {
    min-width: 400px;
    width: auto;
    flex: 0 1 auto;
  }

  &__show-type,
  &__filter-selector {
    width: 200px;
  }

  &__interval {
    width: 300px;
  }

  &__value-filter {
    max-width: 400px;
  }

  &__parameter,
  &__value {
    width: 150px;
  }
}
</style>
