<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.tstop') }}
    v-container
      template(v-if="durationUnit === 'm'")
        v-menu(
        ref="menu"
        :close-on-content-click="false",
        v-model="menu",
        lazy,
        transition="scale-transition",
        )
          v-text-field(slot="activator", :value="dateObject | date('MM/YYYY', true)", readonly)
          v-date-picker(
          ref="picker",
          :value="dateString",
          @input="updateValue",
          :allowed-dates="allowedDates()",
          @change="save"
          )
      template(v-else)
        date-time-picker(@input="updateValue", :value="dateObject", roundHours)
</template>

<script>
import moment from 'moment-timezone';
import { STATS_DURATION_UNITS } from '@/constants';
import DateTimePicker from '@/components/forms/date-time-picker.vue';

export default {
  components: {
    DateTimePicker,
  },
  props: {
    value: {
      type: Number,
      required: true,
    },
    duration: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      menu: false,
    };
  },
  computed: {
    dateObject() {
      return moment.unix(this.value).toDate();
    },
    dateString() {
      return moment.unix(this.value).format('YYYY-MM-DD');
    },
    durationUnit() {
      return this.duration.slice(-1);
    },
  },
  watch: {
    /**
     * Each time the duration change on settings, we need to check if the duration unit is 'm' (for 'Month')
     * If it's 'm', we need to put the date on the 1st day of the month, at 00:00 (UTC)
    */
    durationUnit(value) {
      if (value && value === STATS_DURATION_UNITS.month) {
        const date = moment.tz(this.value * 1000, moment.tz.guess()).startOf('month');
        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      } else {
        this.$emit('input', moment().startOf('hour').unix());
      }
    },
    menu(value) {
      return value && this.$nextTick(() => (this.$refs.picker.activePicker = 'YEAR'));
    },
  },
  methods: {
    updateValue(value) {
      /**
       * If the duration is 'm' (for 'Month'), we need to put the date on the first day of the month,
       * at 00:00 (UTC)
       */
      if (this.durationUnit === STATS_DURATION_UNITS.month) {
        // Get the value's date object, and put it on the first day of the month, at 00:00 (local date)
        const date = moment.tz(value, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      } else {
        this.$emit('input', moment(value).unix());
      }
    },
    save(date) {
      this.$refs.menu.save(date);
    },
    /**
     * Function used by Vuetify date-picker to determine which dates are allowed.
     * Returns a function that filter values
     */
    allowedDates() {
      /**
       * Values are type 'YYYY-MM-DD'
       * We keep only the 'DD' part (with v.split()), parse it to an Int, and keep it only if it's 1
       * Result -> The only available day is the first day of the month
       */
      return v => parseInt(v.split('-')[2], 10) === 1;
    },
  },
};
</script>
