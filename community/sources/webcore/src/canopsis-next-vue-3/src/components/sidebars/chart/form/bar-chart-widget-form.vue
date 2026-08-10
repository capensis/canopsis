<template>
  <c-layout column>
    <field-title v-model="localForm.title" :required="requiredTitle" autofocus />
    <field-periodic-refresh
      v-if="withPeriodicRefresh"
      v-model="localForm.parameters.periodic_refresh"
    />
    <widget-settings-group :title="$t('settings.chart.metricsDisplay')">
      <field-alarm-metric-presets
        v-model="localForm.parameters.metrics"
        :only-external="onlyExternal"
        with-color
        with-external
        only-group
      />
      <field-bar-graph-type v-model="localForm.parameters.stacked" />
    </widget-settings-group>
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-title
        v-model="localForm.parameters.chart_title"
        :label="$t('common.header')"
        :placeholder="$t('settings.headerTitle')"
        name="chart_title"
      />
      <field-quick-date-interval-type
        v-model="localForm.parameters.default_time_range"
        :ranges="metricsRanges"
      />
      <field-sampling v-model="localForm.parameters.default_sampling" />
      <field-filters
        v-if="withFilters"
        v-model="localForm.parameters.mainFilter"
        :filters="localForm.filters"
        addable
        editable
        with-entity
        @update:filters="localForm.filters = $event"
      />
      <field-switcher
        v-model="localForm.parameters.comparison"
        :title="$t('settings.chart.showComparison')"
      />
    </widget-settings-group>
  </c-layout>
</template>

<script>
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import { METRICS_QUICK_RANGES } from '@/constants';

import FieldAlarmMetricPresets from './fields/alarm-metric-presets.vue';
import FieldBarGraphType from './fields/bar-graph-type.vue';
import FieldQuickDateIntervalType from './fields/quick-date-interval-type.vue';
import FieldSampling from './fields/sampling.vue';

export default {
  components: {
    FieldFilters,
    FieldPeriodicRefresh,
    FieldSwitcher,
    FieldTitle,
    WidgetSettingsGroup,
    FieldAlarmMetricPresets,
    FieldBarGraphType,
    FieldQuickDateIntervalType,
    FieldSampling,
  },
  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    withPeriodicRefresh: {
      type: Boolean,
      default: false,
    },
    withFilters: {
      type: Boolean,
      default: false,
    },
    requiredTitle: {
      type: Boolean,
      default: false,
    },
    onlyExternal: {
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
  methods: {
    updateParameter(name, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        parameters: { ...this.modelValue.parameters, [name]: value },
      });
    },
  },
};
</script>
