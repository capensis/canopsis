<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    v-divider
    field-preset(v-model="form.parameters", :type="form.type")
    v-divider
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(v-model="form.parameters.metrics", with-color)
      v-divider
      field-pie-show-mode(v-model="form.parameters.show_mode")
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-chart-title(v-model="form.parameters.chart_title")
      v-divider
      field-quick-date-interval-type(v-model="form.parameters.default_time_range")
      v-divider
      field-sampling(v-model="form.parameters.default_sampling")
      v-divider
      field-alarm-metric-aggregate-function(v-model="form.parameters.aggregate_func")
      v-divider
      field-filters(
        :filters.sync="form.filters",
        addable,
        editable,
        with-alarm,
        with-entity,
        with-pbehavior,
        hide-selector
      )
    v-divider
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';
import FieldTitle from './fields/common/title.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldPreset from './fields/chart/preset.vue';
import FieldAlarmMetricPresets from './fields/chart/alarm-metric-presets.vue';
import FieldPieShowMode from './fields/chart/pie-show-mode.vue';
import FieldChartTitle from './fields/chart/chart-title.vue';
import FieldQuickDateIntervalType from './fields/common/quick-date-interval-type.vue';
import FieldSampling from './fields/common/sampling.vue';
import FieldAlarmMetricAggregateFunction from './fields/chart/alarm-metric-aggregate-function.vue';
import FieldFilters from './fields/common/filters.vue';

export default {
  name: SIDE_BARS.pieChartSettings,
  components: {
    WidgetSettings,
    WidgetSettingsGroup,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldPreset,
    FieldAlarmMetricPresets,
    FieldPieShowMode,
    FieldChartTitle,
    FieldQuickDateIntervalType,
    FieldSampling,
    FieldAlarmMetricAggregateFunction,
    FieldFilters,
  },
  mixins: [widgetSettingsMixin],
};
</script>
