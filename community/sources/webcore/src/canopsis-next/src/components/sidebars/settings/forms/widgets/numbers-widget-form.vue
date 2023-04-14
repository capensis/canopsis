<template lang="pug">
  v-layout(column)
    field-title(v-field="form.title")
    field-periodic-refresh(v-if="withPeriodicRefresh", v-field="form.parameters.periodic_refresh")
    field-preset(v-field="form.parameters", :type="form.type")
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(
        v-field="form.parameters.metrics",
        :parameters="availableParameters",
        with-aggregate-function
      )
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-font-size(v-field="form.parameters.font_size")
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
      field-switcher(v-field="form.parameters.show_trend", :title="$t('settings.chart.showTrend')")
</template>

<script>
import { formMixin } from '@/mixins/form';

import WidgetSettingsGroup from '../../partials/widget-settings-group.vue';
import FieldTitle from '../../fields/common/title.vue';
import FieldPeriodicRefresh from '../../fields/common/periodic-refresh.vue';
import FieldPreset from '../../fields/chart/preset.vue';
import FieldAlarmMetricPresets from '../../fields/chart/alarm-metric-presets.vue';
import FieldBarGraphType from '../../fields/chart/bar-graph-type.vue';
import FieldChartTitle from '../../fields/chart/chart-title.vue';
import FieldQuickDateIntervalType from '../../fields/common/quick-date-interval-type.vue';
import FieldSampling from '../../fields/common/sampling.vue';
import FieldFilters from '../../fields/common/filters.vue';
import FieldSwitcher from '../../fields/common/switcher.vue';
import FieldFontSize from '../../fields/chart/font-size.vue';
import { ALARM_METRIC_PARAMETERS } from '@/constants';

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
    FieldFontSize,
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
    withFilters: {
      type: Boolean,
      default: false,
    },
  },
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
        ALARM_METRIC_PARAMETERS.timeToAck,
        ALARM_METRIC_PARAMETERS.timeToResolve,
        ALARM_METRIC_PARAMETERS.notAckedAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInHourAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms,
        ALARM_METRIC_PARAMETERS.notAckedInDayAlarms,
      ];
    },
  },
  methods: {
    updateFilters(filters) {
      this.updateField('filters', filters);
    },
  },
};
</script>
