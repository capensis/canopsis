<template>
  <v-layout class="availability-widget-filters col-gap-6 row-gap-3" wrap>
    <c-availability-advanced-search
      :searches="searches"
      class="mt-0 availability-widget-filters__search"
      @submit="updateSearch"
      @reset="resetSearch"
      @toggle-pin:search="togglePinSearchInUserPreferences"
      @remove:search="removeSearchFromUserPreferences"
    />
    <c-quick-date-interval-field
      v-if="showInterval"
      :interval="interval"
      :quick-ranges="quickRanges"
      class="availability-widget-filters__interval"
      short
      with-hours
      allow-future
      @input="$emit('update:interval', $event)"
    />
    <v-layout
      v-if="showFilter"
      class="availability-widget-filters__filter-wrapper"
      align-end
    >
      <filter-selector
        :label="$t('settings.selectAFilter')"
        :filters="userFilters"
        :locked-filters="widgetFilters"
        :locked-value="lockedFilter"
        :value="filters"
        :disabled="filterDisabled"
        class="availability-widget-filters__filter-selector mr-4"
        clearable
        hide-details
        @input="$emit('update:filters', $event)"
      />
      <filters-list-btn
        :widget-id="widgetId"
        addable
        editable
        private
        with-entity
      />
    </v-layout>

    <availability-display-parameter-field
      :value="displayParameter"
      :label="$t('common.value')"
      class="availability-widget-filters__parameter"
      @input="$emit('update:display-parameter', $event)"
    />
    <availability-show-type-field
      :value="type"
      :label="$t('common.show')"
      class="availability-widget-filters__show-type"
      @input="$emit('update:type', $event)"
    />
    <c-enabled-field
      :value="trend"
      :label="$t('settings.showTrend')"
      hide-details
      @input="$emit('update:trend', $event)"
    />
    <availability-value-filter-field
      v-model="localValueFilter"
      :show-type="type"
      :max-seconds="maxValueFilterSeconds"
      class="availability-widget-filters__value-filter"
      @input="handleUpdateValueFilter"
    />
    <c-action-btn
      v-if="showExport"
      :loading="exporting"
      :tooltip="$t('settings.exportAsCsv')"
      icon="cloud_download"
      @click="$emit('export')"
    />
  </v-layout>
</template>

<script>
import { debounce } from 'lodash';
import { ref, toRef, computed, watch } from 'vue';

import { AVAILABILITY_QUICK_RANGES } from '@/constants';

import { useWidgetAdvancedSearchSavedItems } from '@/hooks/widget/advanced-search-saved-items';

import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';
import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import AvailabilityDisplayParameterField from '@/components/other/availability/form/fields/availability-display-parameter-field.vue';
import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';
import AvailabilityValueFilterField from '@/components/other/availability/form/fields/availability-value-filter-field.vue';

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
    columns: {
      type: Array,
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

    const {
      searches,
      updateSearch,
      resetSearch,
      togglePinSearchInUserPreferences,
      removeSearchFromUserPreferences,
    } = useWidgetAdvancedSearchSavedItems(
      {
        widgetId: toRef(props, 'widgetId'),
        query: toRef(props, 'query'),
      },
      emit,
    );

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

      searches,
      updateSearch,
      resetSearch,
      togglePinSearchInUserPreferences,
      removeSearchFromUserPreferences,
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
  }

  &__show-type, &__filter-selector {
    width: 200px;
  }

  &__interval {
    width: 300px;
  }

  &__value-filter {
    max-width: 400px;
  }

  &__parameter, &__value {
    width: 150px;
  }
}
</style>
