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
              :date-object-preparer="startDateObjectPreparer",
              :round-hours="roundHours",
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
              :date-object-preparer="stopDateObjectPreparer",
              :round-hours="roundHours",
              name="tstop",
              @update:objectValue="$emit('update:stopObjectValue', $event)"
            )
      v-flex.pl-1(xs6, data-test="intervalRange")
        v-select(
          v-model="range",
          :items="quickRanges",
          :label="$t('quickRanges.title')",
          return-object
        )
</template>

<script>
import { TIME_UNITS, QUICK_RANGES, DATETIME_FORMATS } from '@/constants';

import { prepareDateToObject, findQuickRangeValue } from '@/helpers/date/date-intervals';
import { convertDateToStartOfUnitString, subtractUnitFromDate } from '@/helpers/date/date';

import { formMixin } from '@/mixins/form';

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
        const range = findQuickRangeValue(tstart, tstop);

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
              periodUnit: TIME_UNITS.hour,
              periodValue: 1,

              tstart: convertDateToStartOfUnitString(
                subtractUnitFromDate(Date.now(), 1, TIME_UNITS.hour),
                TIME_UNITS.hour,
                DATETIME_FORMATS.dateTimePicker,
              ),

              tstop: convertDateToStartOfUnitString(
                Date.now(),
                TIME_UNITS.hour,
                DATETIME_FORMATS.dateTimePicker,
              ),
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
      return Object.values(QUICK_RANGES).map(range => ({
        ...range,

        text: this.$t(`quickRanges.types.${range.value}`),
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
