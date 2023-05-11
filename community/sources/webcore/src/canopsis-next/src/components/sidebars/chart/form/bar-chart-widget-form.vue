<template lang="pug">
  v-layout(column)
    field-title(v-field="form.title", :required="requiredTitle")
    field-periodic-refresh(v-if="withPeriodicRefresh", v-field="form.parameters.periodic_refresh")
    field-preset(v-if="withPreset", v-field="form.parameters", :type="form.type")
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(
        v-field="form.parameters.metrics",
        :only-external="onlyExternal",
        with-color,
        with-external,
        only-group
      )
      field-bar-graph-type(v-field="form.parameters.stacked")
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-chart-title(v-field="form.parameters.chart_title")
      field-quick-date-interval-type(v-field="form.parameters.default_time_range")
      field-sampling(v-field="form.parameters.default_sampling")
      field-filters(
        v-if="withFilters",
        :filters="form.filters",
        addable,
        editable,
        with-entity,
        hide-selector,
        @update:filters="updateFilters"
      )
      field-switcher(v-field="form.parameters.comparison", :title="$t('settings.chart.showComparison')")
</template>

<script>
import { formMixin } from '@/mixins/form';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldPreset from '@/components/sidebars/chart/form/fields/preset.vue';
import FieldAlarmMetricPresets from '@/components/sidebars/chart/form/fields/alarm-metric-presets.vue';
import FieldBarGraphType from '@/components/sidebars/chart/form/fields/bar-graph-type.vue';
import FieldChartTitle from '@/components/sidebars/chart/form/fields/chart-title.vue';
import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/quick-date-interval-type.vue';
import FieldSampling from '@/components/sidebars/chart/form/fields/sampling.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';

export default {
  components: {
    WidgetSettingsGroup,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldPreset,
    FieldAlarmMetricPresets,
    FieldBarGraphType,
    FieldChartTitle,
    FieldQuickDateIntervalType,
    FieldSampling,
    FieldFilters,
    FieldSwitcher,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    withPeriodicRefresh: {
      type: Boolean,
      default: false,
    },
    withPreset: {
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
  methods: {
    updateFilters(filters) {
      this.updateField('filters', filters);
    },
  },
};
</script>
