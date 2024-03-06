<template>
  <v-layout
    class="availability-widget-filters gap-6"
    align-end
    wrap
  >
    <c-quick-date-interval-field
      v-if="showInterval"
      :interval="interval"
      :min="minIntervalDate"
      :quick-ranges="quickRanges"
      class="availability-widget-filters__interval"
      short
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
      :value="valueFilter"
      :show-type="type"
      :max-value="intervalSecondsDiff"
      class="availability-widget-filters__value-filter"
      @input="debouncedUpdateValueFilter"
    />
    <c-action-btn
      :loading="exporting"
      :tooltip="$t('settings.exportAsCsv')"
      icon="cloud_download"
      @click="$emit('export')"
    />
  </v-layout>
</template>

<script>
import { debounce } from 'lodash';

import { AVAILABILITY_QUICK_RANGES } from '@/constants';

import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';
import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import AvailabilityDisplayParameterField from '@/components/other/availability/form/fields/availability-display-parameter-field.vue';
import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';
import AvailabilityValueFilterField
  from '@/components/widgets/availability/form/fields/availability-value-filter-field.vue';

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
    minIntervalDate: {
      type: Number,
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
    intervalSecondsDiff: {
      type: Number,
      required: false,
    },
  },
  computed: {
    quickRanges() {
      return Object.values(AVAILABILITY_QUICK_RANGES);
    },
  },
  created() {
    this.debouncedUpdateValueFilter = debounce(this.updateValueFilter, 1000);
  },
  methods: {
    updateValueFilter(value) {
      this.$emit('update:value-filter', value);
    },
  },
};
</script>

<style lang="scss" scoped>
.availability-widget-filters {
  &__filter-wrapper {
    flex-grow: 0;
  }

  &__show-type, &__filter-selector, &__interval {
    max-width: 200px;
  }

  &__value-filter {
    max-width: 400px;
  }

  &__parameter, &__value {
    max-width: 150px;
  }
}
</style>
