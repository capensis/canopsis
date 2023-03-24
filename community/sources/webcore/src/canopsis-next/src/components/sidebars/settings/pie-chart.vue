<template lang="pug">
  widget-settings(:submitting="submitting", divider, @submit="submit")
    field-title(v-model="form.title")
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    field-preset(v-model="form.parameters", :type="form.type")
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(
        v-model="form.parameters.metrics",
        :parameters="availableParameters",
        :min="2",
        with-color
      )
      field-pie-show-mode(v-model="form.parameters.show_mode")
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-chart-title(v-model="form.parameters.chart_title")
      field-quick-date-interval-type(v-model="form.parameters.default_time_range")
      field-sampling(v-model="form.parameters.default_sampling")
      field-alarm-metric-aggregate-function(v-model="form.parameters.aggregate_func")
      field-filters(
        :filters.sync="form.filters",
        addable,
        editable,
        with-alarm,
        with-entity,
        with-pbehavior,
        hide-selector
      )
</template>

<script>
import { ALARM_METRIC_PARAMETERS, SIDE_BARS } from '@/constants';

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
  computed: {
    availableParameters() {
      return [
        ALARM_METRIC_PARAMETERS.createdAlarms,
        ALARM_METRIC_PARAMETERS.activeAlarms,
        ALARM_METRIC_PARAMETERS.nonDisplayedAlarms,
        ALARM_METRIC_PARAMETERS.instructionAlarms,
        ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms,
        ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms,
        ALARM_METRIC_PARAMETERS.pbehaviorAlarms,
        ALARM_METRIC_PARAMETERS.correlationAlarms,
        ALARM_METRIC_PARAMETERS.ackAlarms,
        ALARM_METRIC_PARAMETERS.cancelAckAlarms,
        ALARM_METRIC_PARAMETERS.ackActiveAlarms,
        ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
        ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms,
        ALARM_METRIC_PARAMETERS.notAckedAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInHourAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInDayAlarms,
      ];
    },
  },
};
</script>
