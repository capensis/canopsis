<template lang="pug">
  v-layout(column)
    field-title(v-field="form.title", :required="requiredTitle")
    field-periodic-refresh(v-if="withPeriodicRefresh", v-field="form.parameters")
    field-preset(v-if="withPreset", v-field="form.parameters", :type="form.type")
    widget-settings-group(:title="$t('settings.chart.metricsDisplay')")
      field-alarm-metric-presets(
        v-field="form.parameters.metrics",
        :parameters="availableParameters",
        :min="2",
        with-color
      )
      field-pie-show-mode(v-field="form.parameters.show_mode")
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-title(
        v-field="form.parameters.chart_title",
        :label="$tc('common.header')",
        :placeholder="$t('settings.headerTitle')",
        name="chart_title"
      )
      field-quick-date-interval-type(v-field="form.parameters.default_time_range")
      field-sampling(v-field="form.parameters.default_sampling")
      field-alarm-metric-aggregate-function(v-field="form.parameters.aggregate_func")
      field-filters(
        v-if="withFilters",
        v-field="form.parameters.mainFilter",
        :filters="form.filters",
        addable,
        editable,
        with-entity,
        @update:filters="updateFilters"
      )
</template>

<script>
import { ALARM_METRIC_PARAMETERS } from '@/constants';

import { formMixin } from '@/mixins/form';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldQuickDateIntervalType from '@/components/sidebars/chart/form/fields/quick-date-interval-type.vue';
import FieldPreset from '@/components/sidebars/chart/form/fields/preset.vue';
import FieldAlarmMetricPresets from '@/components/sidebars/chart/form/fields/alarm-metric-presets.vue';
import FieldPieShowMode from '@/components/sidebars/chart/form/fields/pie-show-mode.vue';
import FieldSampling from '@/components/sidebars/chart/form/fields/sampling.vue';
import FieldAlarmMetricAggregateFunction from '@/components/sidebars/chart/form/fields/alarm-metric-aggregate-function.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';

export default {
  components: {
    WidgetSettingsGroup,
    FieldTitle,
    FieldPeriodicRefresh,
    FieldQuickDateIntervalType,
    FieldPreset,
    FieldAlarmMetricPresets,
    FieldPieShowMode,
    FieldSampling,
    FieldAlarmMetricAggregateFunction,
    FieldFilters,
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
