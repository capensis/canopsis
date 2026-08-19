<template>
  <c-layout column>
    <field-title v-model="localForm.title" :required="requiredTitle" autofocus />
    <field-periodic-refresh
      v-if="withPeriodicRefresh"
      v-model="localForm.parameters.periodic_refresh"
    />
    <field-preset
      v-if="withPreset"
      v-model="localForm.parameters"
      :type="localForm.type"
    />
    <widget-settings-group :title="$t('settings.chart.metricsDisplay')">
      <field-alarm-metric-presets
        v-model="localForm.parameters.metrics"
        :parameters="availableParameters"
        :min="2"
        with-color
      />
      <field-pie-show-mode v-model="localForm.parameters.show_mode" />
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
      <field-alarm-metric-aggregate-function v-model="localForm.parameters.aggregate_func" />
      <field-filters
        v-if="withFilters"
        v-model="localForm.parameters.mainFilter"
        :filters="localForm.filters"
        addable
        editable
        with-entity
        @update:filters="localForm.filters = $event"
      />
    </widget-settings-group>
  </c-layout>
</template>

<script>
import { ALARM_METRIC_PARAMETERS, METRICS_QUICK_RANGES } from '@/constants';

import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

import FieldAlarmMetricAggregateFunction from './fields/alarm-metric-aggregate-function.vue';
import FieldAlarmMetricPresets from './fields/alarm-metric-presets.vue';
import FieldPieShowMode from './fields/pie-show-mode.vue';
import FieldPreset from './fields/preset.vue';
import FieldQuickDateIntervalType from './fields/quick-date-interval-type.vue';
import FieldSampling from './fields/sampling.vue';

export default {
  components: {
    FieldAlarmMetricAggregateFunction,
    FieldAlarmMetricPresets,
    FieldFilters,
    FieldPeriodicRefresh,
    FieldPieShowMode,
    FieldPreset,
    FieldQuickDateIntervalType,
    FieldSampling,
    FieldTitle,
    WidgetSettingsGroup,
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
