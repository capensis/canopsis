<template>
  <c-layout column>
    <field-title v-model="localForm.title" />
    <field-periodic-refresh v-model="localForm.parameters.periodic_refresh" />
    <field-main-parameter
      v-model="localForm.parameters.mainParameter"
      :type="type"
      :widget-id="widget._id"
    />
    <field-statistics-columns v-model="localForm.parameters.widgetColumns" :type="type" />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-title
        v-model="localForm.parameters.table_title"
        :label="$t('common.header')"
        :placeholder="$t('settings.headerTitle')"
        name="table_title"
      />
      <field-quick-date-interval-type
        v-model="localForm.parameters.default_time_range"
        :ranges="metricsRanges"
      />
      <field-filters
        v-if="showFilter"
        v-model="localForm.parameters.mainFilter"
        :filters="localForm.filters"
        :widget-id="widget._id"
        addable
        editable
        with-entity
        @update:filters="localForm.filters = $event"
      />
    </widget-settings-group>
  </c-layout>
</template>

<script>
import { KPI_RATING_SETTINGS_TYPES, METRICS_QUICK_RANGES } from '@/constants';

import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/quick-date-interval-type.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

import FieldMainParameter from './fields/main-parameter.vue';
import FieldStatisticsColumns from './fields/statistics-columns.vue';

export default {
  components: {
    WidgetSettingsGroup,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldQuickDateIntervalType,
    FieldFilters,
    FieldMainParameter,
    FieldStatisticsColumns,
  },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
    showFilter: {
      type: Boolean,
      default: false,
    },
    showInterval: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    localForm: {
      get() {
        return this.modelValue;
      },
      set(value) {
        this.$emit('update:modelValue', value);
      },
    },

    metricsRanges() {
      return METRICS_QUICK_RANGES;
    },
  },
};
</script>
