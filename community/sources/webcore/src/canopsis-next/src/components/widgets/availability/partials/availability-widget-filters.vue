<template>
  <div class="availability-widget-filters col-gap-6 row-gap-3">
    <c-advanced-search-field
      :query.sync="searchQuery"
      :tooltip="$t('context.advancedSearch')"
      :columns="columns"
      class="pa-0 availability-widget-filters__search"
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
        v-if="filterAddable || filterEditable"
        :widget-id="widgetId"
        :addable="filterAddable"
        :editable="filterEditable"
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
  </div>
</template>

<script>
import { debounce } from 'lodash';
import { ref, watch, computed } from 'vue';

import { AVAILABILITY_QUICK_RANGES } from '@/constants';

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
    search: {
      type: String,
      required: false,
    },
    columns: {
      type: Array,
      required: false,
    },
    interval: {
      type: Object,
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
    showInterval: {
      type: Boolean,
      default: false,
    },
    showFilter: {
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
    displayParameter: {
      type: Number,
      required: false,
    },
    type: {
      type: Number,
      required: false,
    },
    trend: {
      type: Boolean,
      required: false,
    },
    exporting: {
      type: Boolean,
      default: false,
    },
    valueFilter: {
      type: Object,
      required: false,
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

    const searchQuery = computed({
      get() {
        return {
          search: props.search,
        };
      },
      set({ search }) {
        emit('update:search', search);
      },
    });

    watch(
      () => props.valueFilter,
      () => {
        localValueFilter.value = props.valueFilter && { ...props.valueFilter };
      },
      { immediate: true },
    );

    const emitUpdateValueFilter = valueFilter => emit('update:value-filter', valueFilter);
    const debouncedEmitUpdateValueFilter = debounce(emitUpdateValueFilter, 1000);

    const handleUpdateValueFilter = (valueFilter) => {
      if (!valueFilter || valueFilter.value === props.valueFilter?.value) {
        emitUpdateValueFilter(valueFilter);
      } else {
        debouncedEmitUpdateValueFilter(valueFilter);
      }
    };

    return {
      searchQuery,
      localValueFilter,
      quickRanges,

      handleUpdateValueFilter,
    };
  },
};
</script>

<style lang="scss" scoped>
.availability-widget-filters {
  display: flex;
  flex-wrap: wrap;
  align-items: end;

  &__search {
    width: 400px;
  }

  &__filter-wrapper {
    flex-grow: 0;
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
