<template lang="pug">
  div
    v-layout
      v-flex(xs6)
        v-layout(align-center)
          v-flex
            date-time-picker-text-field(
              data-test="intervalStart",
              v-model="tstartDateString",
              v-validate="tstartRules",
              :label="$t('common.startDate')",
              :dateObjectPreparer="startDateObjectPreparer",
              :roundHours="roundHours",
              name="tstart",
              @update:objectValue="$emit('update:startObjectValue', $event)"
            )
        v-layout(align-center)
          v-flex
            date-time-picker-text-field(
              data-test="intervalStop",
              v-model="tstopDateString",
              v-validate="tstopRules",
              :label="$t('common.endDate')",
              :dateObjectPreparer="stopDateObjectPreparer",
              :roundHours="roundHours",
              name="tstop",
              @update:objectValue="$emit('update:stopObjectValue', $event)"
            )
      v-flex.pl-1(xs6, data-test="intervalRange")
        v-select(
          v-model="range",
          :items="quickRanges",
          :label="$t('settings.statsDateInterval.fields.quickRanges')",
          return-object
        )
</template>

<script>
import moment from 'moment';

import { STATS_DURATION_UNITS, STATS_QUICK_RANGES, DATETIME_FORMATS } from '@/constants';

import { prepareDateToObject, findRange } from '@/helpers/date-intervals';

import formMixin from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerTextField,
  },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    roundHours: {
      type: Boolean,
      default: false,
    },
    value: {
      type: Object,
      required: true,
    },
    tstopRules: {
      type: [String, Array],
      default: null,
    },
    tstartRules: {
      type: [String, Array],
      default: null,
    },
  },
  computed: {
    stopDateObjectPreparer() {
      return this.preparerDateToObjectGetter('stop');
    },
    startDateObjectPreparer() {
      return this.preparerDateToObjectGetter('start');
    },
    range: {
      get() {
        const { tstart, tstop } = this.value;
        const range = findRange(tstart, tstop);

        return this.quickRanges.find(({ value }) => value === range.value);
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
    preparerDateToObjectGetter(type) {
      return date => prepareDateToObject(date, type, this.roundHours ? 'hour' : 'minute');
    },
  },
};
</script>
