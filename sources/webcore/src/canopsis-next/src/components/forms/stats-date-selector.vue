<template lang="pug">
  div
    v-layout
      v-flex(xs5)
        v-layout(align-center)
          v-text-field(label="Start date", v-model="tstartDateString")
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
            date-time-picker(:value="toDateObject(value.tstart)", @submit="handleTstartChange", roundHours)
        v-layout(align-center)
          v-text-field(label="End date", v-model="tstopDateString")
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
            date-time-picker(:value="toDateObject(value.tstop)", @submit="handleTstopChange", roundHours)
      v-flex.px-1(xs7)
        h4.primary--text Quick ranges
        v-layout(wrap)
          v-flex(xs6, v-for="range in $constants.STATS_QUICK_RANGES", :key="range.value")
            a.black--text(
            @click="updateRange(range.start, range.stop)",
            ) - {{ $t(`settings.statsDateInterval.quickRanges.${range.value}`) }}
</template>

<script>
import moment from 'moment';

import formMixin from '@/mixins/form';

import DateTimePicker from '@/components/forms/fields/date-picker/date-time-picker-field.vue';

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
    tstartDateString: {
      get() {
        return moment(this.value.tstart).isValid() ? moment(this.value.tstart).format('DD/MM/YYYY HH:mm') : this.value.tstart;
      },
      set(value) {
        if (value !== this.value.tstart) {
          this.updateField('tstart', value);
        }
      },
    },
    tstopDateString: {
      get() {
        return moment(this.value.tstop).isValid() ? moment(this.value.tstop).format('DD/MM/YYYY HH:mm') : this.value.tstop;
      },
      set(value) {
        if (value !== this.value.tstop) {
          this.updateField('tstop', value);
        }
      },
    },
  },
  methods: {
    toDateObject(date) {
      const isValidDate = moment(date).isValid();

      if (isValidDate) {
        return this.periodUnit === 'm' ?
          moment(date).startOf('month').toDate() :
          moment(date).startOf('hour').toDate();
      }

      return this.periodUnit === 'm' ?
        moment().startOf('month').toDate() :
        moment().startOf('hour').toDate();
    },

    handleTstartChange(tstart) {
      this.updateField('tstart', tstart);
      this.isTstartDateMenuOpen = false;
    },

    handleTstopChange(tstop) {
      this.updateField('tstop', tstop);
      this.isTstopDateMenuOpen = false;
    },

    updateRange(tstart, tstop) {
      this.$emit('input', { tstart, tstop });
    },
  },
};
</script>
