<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.tstop') }}
    v-container
      template(v-if="durationUnit === 'm'")
        v-menu(
        ref="menu"
        :close-on-content-click="false",
        transition="scale-transition",
        )
          v-text-field(slot="activator", :value="dateObject | date('MM/YYYY', true)", readonly)
          v-date-picker(@input="updateValue", :allowed-dates="allowedDates")
      template(v-else)
        date-time-picker(:value="dateObject", @input="updateValue")
</template>

<script>
import moment from 'moment-timezone';
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
  computed: {
    dateObject() {
      return moment.unix(this.value).toDate();
    },
    durationUnit() {
      return this.duration.slice(-1);
    },
  },
  methods: {
    updateValue(value) {
      if (this.durationUnit === 'm') {
        const date = moment.tz(value, moment.tz.guess()).startOf('month');

        this.$emit('input', date.add(date.utcOffset(), 'm').unix());
      } else {
        this.$emit('input', moment(value).unix());
      }
    },
    allowedDates: val => parseInt(val.split('-')[2], 10) === 1,
  },
};
</script>
