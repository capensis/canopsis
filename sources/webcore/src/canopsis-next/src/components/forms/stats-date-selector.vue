<template lang="pug">
  div
    v-layout
      v-flex(xs6)
        v-layout(align-center)
          v-text-field(v-model="tstartDateString", :label="$t('common.startDate')")
          v-menu(
          ref="menu",
          v-model="isTstartDateMenuOpen",
          content-class="date-time-picker",
          transition="slide-y-transition",
          max-width="290px",
          :close-on-content-click="false",
          right,
          lazy,
          )
            v-btn(slot="activator", icon, fab, small, color="secondary")
              v-icon calendar_today
            date-time-picker(
            :value="toDateObject(value.tstart)",
            :opened="isTstartDateMenuOpen",
            roundHours,
            @input="updateTstartField"
            )
        v-layout(align-center)
          v-text-field(v-model="tstopDateString", :label="$t('common.endDate')")
          v-menu(
          v-model="isTstopDateMenuOpen",
          content-class="date-time-picker",
          transition="slide-y-transition",
          max-width="290px",
          :close-on-content-click="false",
          right,
          lazy,
          )
            v-btn(slot="activator", icon, fab, small, color="secondary")
              v-icon calendar_today
            date-time-picker(
            :value="toDateObject(value.tstop)",
            :opened="isTstopDateMenuOpen",
            @input="updateTstopField",
            roundHours
            )
      v-flex.px-1(xs6)
        v-select(v-model="range", :items="quickRanges", label="Quick ranges", return-object)
</template>

<script>
import moment from 'moment';

import { STATS_DURATION_UNITS, STATS_QUICK_RANGES, DATETIME_FORMATS } from '@/constants';

import formMixin from '@/mixins/form';

import DateTimePicker from '@/components/forms/fields/date-time-picker/date-time-picker.vue';

export default {
  components: {
    DateTimePicker,
  },
  mixins: [formMixin],
  props: {
    periodUnit: {
      type: String,
      required: true,
    },
    value: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isTstartDateMenuOpen: false,
      isTstopDateMenuOpen: false,
    };
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
          let tstart = range.start;
          let tstop = range.stop;

          if (!tstop || !tstart) {
            tstart = moment().startOf('h').format(DATETIME_FORMATS.dateTimePicker);
            tstop = moment().add(1, 'h').startOf('h').format(DATETIME_FORMATS.dateTimePicker);
          }

          this.$emit('input', { tstart, tstop });
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

    toDateObject() {
      return (date) => {
        const unit = this.periodUnit === STATS_DURATION_UNITS.month ? 'month' : 'hour';
        const momentDate = moment(date, DATETIME_FORMATS.dateTimePicker);

        if (momentDate.isValid()) {
          return momentDate.startOf(unit).toDate();
        }

        return moment().startOf(unit).toDate();
      };
    },
  },
  methods: {
    updateTstartField(tstart) {
      this.updateField('tstart', moment(tstart).format(DATETIME_FORMATS.dateTimePicker));
    },

    updateTstopField(tstop) {
      this.updateField('tstop', moment(tstop).format(DATETIME_FORMATS.dateTimePicker));
    },
  },
};
</script>
