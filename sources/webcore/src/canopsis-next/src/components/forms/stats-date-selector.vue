<template lang="pug">
  div
    v-layout
      v-flex(xs6)
        v-layout(align-center)
          date-time-picker-text-field(
          ref="tstart",
          v-model="tstartDateString",
          :label="$t('common.startDate')",
          :dateObjectPreparer="getDateObjectPreparer('start')",
          name="tstart",
          roundHours
          )
        v-layout(align-center)
          date-time-picker-text-field(
          v-model="tstopDateString",
          v-validate="'after_custom:tstart'",
          :label="$t('common.endDate')",
          :dateObjectPreparer="getDateObjectPreparer('stop')",
          name="tstop",
          roundHours
          )
      v-flex.px-1(xs6)
        v-select(v-model="range", :items="quickRanges", label="Quick ranges", return-object)
</template>

<script>
import moment from 'moment';

import { STATS_DURATION_UNITS, STATS_QUICK_RANGES, DATETIME_FORMATS } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

import formMixin from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerTextField,
  },
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  computed: {
    range: {
      get() {
        const activeRange = this.quickRanges
          .find(range => this.value.tstart === range.start && this.value.tstop === range.stop);

        if (!activeRange) {
          return this.quickRanges.find(range => range.value === STATS_QUICK_RANGES.custom.value);
        }

        return activeRange;
      },
      set(range) {
        if (range.value !== this.range.value) {
          let newValue = {
            tstart: range.start,
            tstop: range.stop,
          };

          if (!newValue.tstop || !newValue.tstart) {
            newValue = {
              periodUnit: STATS_DURATION_UNITS.hour,
              periodValue: 1,

              tstart: moment()
                .subtract(1, STATS_DURATION_UNITS.hour)
                .startOf(STATS_DURATION_UNITS.hour)
                .format(DATETIME_FORMATS.dateTimePicker),

              tstop: moment()
                .startOf(STATS_DURATION_UNITS.hour)
                .format(DATETIME_FORMATS.dateTimePicker),
            };
          }

          this.updateModel({
            ...this.value,
            ...newValue,
          });
        }
      },
    },

    quickRanges() {
      return Object.values(STATS_QUICK_RANGES).map(range => ({
        ...range,

        text: this.$t(`settings.statsDateInterval.quickRanges.${range.value}`),
      }));
    },

    tstartDateString: {
      get() {
        return this.value.tstart;
      },
      set(value) {
        if (value !== this.value.tstart) {
          this.updateField('tstart', value);
        }
      },
    },

    tstopDateString: {
      get() {
        return this.value.tstop;
      },
      set(value) {
        if (value !== this.value.tstop) {
          this.updateField('tstop', value);
        }
      },
    },
  },
  methods: {
    toDateObject(date, type) {
      const unit = this.periodUnit === STATS_DURATION_UNITS.month ? 'month' : 'hour';
      const momentDate = dateParse(date, type, DATETIME_FORMATS.dateTimePicker);

      if (momentDate.isValid()) {
        return momentDate.startOf(unit).toDate();
      }

      return moment().startOf(unit).toDate();
    },
    getDateObjectPreparer(type) {
      return value => this.toDateObject(value, type);
    },
  },
};
</script>
