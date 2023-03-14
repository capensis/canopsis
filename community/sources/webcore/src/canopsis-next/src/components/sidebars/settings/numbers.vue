<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-periodic-refresh(v-model="form.parameters.periodic_refresh")
    v-divider
    field-preset(v-model="form.parameters", :type="form.type")
    v-divider
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(
        v-model="form.parameters.metrics",
        :parameters="availableParameters",
        with-aggregate-function
      )
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-chart-title(v-model="form.parameters.chart_title")
      v-divider
      field-quick-date-interval-type(v-model="form.parameters.default_time_range")
      v-divider
      field-sampling(v-model="form.parameters.default_sampling")
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
      field-switcher(v-model="form.parameters.show_trend", :title="$t('settings.chart.showTrend')")
    v-divider
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
import FieldBarGraphType from './fields/chart/bar-graph-type.vue';
import FieldChartTitle from './fields/chart/chart-title.vue';
import FieldQuickDateIntervalType from './fields/common/quick-date-interval-type.vue';
import FieldSampling from './fields/common/sampling.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldSwitcher from './fields/common/switcher.vue';

export default {
  name: SIDE_BARS.numbersSettings,
  components: {
    WidgetSettings,
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
        ALARM_METRIC_PARAMETERS.ratioCorrelation,
        ALARM_METRIC_PARAMETERS.ratioInstructions,
        ALARM_METRIC_PARAMETERS.ratioTickets,
        ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
        ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms,
        ALARM_METRIC_PARAMETERS.averageAck,
        ALARM_METRIC_PARAMETERS.averageResolve,
        ALARM_METRIC_PARAMETERS.notAckedAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInHourAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInDayAlarms,
      ];
    },
  },
};
</script>
