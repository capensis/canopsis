<template lang="pug">
  div
    v-container
      v-layout
        v-flex(xs3, v-if="!hiddenFields.includes('periodValue')")
          v-text-field.pt-0(
            v-field="form.periodValue",
            :label="$t('modals.statsDateInterval.fields.periodValue')",
            type="number",
            data-test="intervalPeriodValue"
          )
        v-flex(data-test="intervalPeriodUnit")
          v-select.pt-0(
            v-field="form.periodUnit",
            :items="periodUnits",
            :label="$t('modals.statsDateInterval.fields.periodUnit')"
          )
      v-alert.mb-2(
        :value="isPeriodMonth",
        type="info"
      ) {{ $t('settings.statsDateInterval.monthPeriodInfo') }}
      date-interval-selector.my-1(
        v-field="form",
        :tstopRules="tstopRules",
        roundHours,
        @update:startObjectValue="updateStartObjectValue",
        @update:stopObjectValue="updateStopObjectValue"
      )
      v-alert.mb-2(
        :value="isPeriodMonth",
        type="info"
      ) {{ monthIntervalMessage }}
</template>

<script>
import { DATETIME_FORMATS, STATS_DURATION_UNITS } from '@/constants';

import {
  dateParse,
  prepareStatsStopForMonthPeriod,
  prepareStatsStartForMonthPeriod,
} from '@/helpers/date/date-intervals';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

export default {
  inject: ['$validator'],
  components: { DateIntervalSelector },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    hiddenFields: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      dateObjectValues: {
        start: null,
        stop: null,
      },
    };
  },
  computed: {
    tstopRules() {
      return `after_custom:${this.form.tstart}`;
    },

    periodUnits() {
      return [
        {
          text: this.$tc('common.times.hour'),
          value: STATS_DURATION_UNITS.hour,
        },
        {
          text: this.$tc('common.times.day'),
          value: STATS_DURATION_UNITS.day,
        },
        {
          text: this.$tc('common.times.week'),
          value: STATS_DURATION_UNITS.week,
        },
        {
          text: this.$tc('common.times.month'),
          value: STATS_DURATION_UNITS.month,
        },
      ];
    },

    isPeriodMonth() {
      return this.form.periodUnit === STATS_DURATION_UNITS.month;
    },

    monthIntervalMessage() {
      return this.$t('modals.statsDateInterval.info.monthPeriodUnit', {
        start: this.$options.filters.date(this.dateObjectValues.start, 'long', false),
        stop: this.$options.filters.date(this.dateObjectValues.stop, 'long', false),
      });
    },
  },
  created() {
    this.$validator.extend('after_custom', {
      getMessage: () => this.$t('modals.statsDateInterval.errors.endDateLessOrEqualStartDate'),
      validate: (value, [otherValue]) => {
        try {
          const convertedStop = dateParse(value, 'stop', DATETIME_FORMATS.dateTimePicker);
          const convertedStart = dateParse(otherValue, 'start', DATETIME_FORMATS.dateTimePicker);

          return !convertedStop.isSameOrBefore(convertedStart);
        } catch (err) {
          return true; // TODO: problem with i18n: https://github.com/baianat/vee-validate/issues/2025
        }
      },
    }, {
      hasTarget: true,
    });
  },
  methods: {
    updateStartObjectValue(value) {
      this.dateObjectValues.start = value && prepareStatsStartForMonthPeriod(value);
    },

    updateStopObjectValue(value) {
      this.dateObjectValues.stop = value && prepareStatsStopForMonthPeriod(value);
    },
  },
};
</script>
