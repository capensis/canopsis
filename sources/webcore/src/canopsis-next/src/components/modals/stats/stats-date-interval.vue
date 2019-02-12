<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Stats - Date interval
    v-card-text
      v-select.pt-0(
      :items="periodUnits",
      v-model="form.durationUnit",
      label="Period",
      )
      v-card
        v-card-title Start date
        v-card-text
          template(v-if="form.durationUnit === 'm'")
            v-menu(
            ref="menu"
            :close-on-content-click="false",
            v-model="menu",
            transition="scale-transition",
            )
              v-text-field(slot="activator", v-model="form.tstart.format('YYYY-MM-DD')", readonly)
              date-picker(
              ref="picker",
              :value="form.tstart.format('YYYY-MM-DD')",
              @input="",
              @change="save",
              type="month",
              year-first,
              )
          template(v-else)
            date-time-picker(v-model="form.tstart.toDate()", roundHours)
      v-card.my-2
        v-card-title End date
        v-card-text
          v-switch(label="Fixed date", v-model="form.tstopMode")
          template(v-if="form.durationUnit === 'm'")
            v-menu(
            ref="menu"
            :close-on-content-click="false",
            v-model="menu",
            transition="scale-transition",
            )
              v-text-field(slot="activator", v-model="form.tstart.format('YYYY-MM-DD')", readonly)
              date-picker(
              ref="picker",
              :value="form.tstart.format('YYYY-MM-DD')",
              @input="",
              @change="save",
              type="month",
              year-first,
              )
          template(v-else)
            date-time-picker(v-model="form.tstart.toDate()", roundHours)
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="") {{ $t('common.submit') }}
</template>

<script>
import moment from 'moment-timezone';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import DatePicker from '@/components/forms/fields/date-picker/date-picker.vue';
import DateTimePicker from '@/components/forms/fields/date-time-picker.vue';

export default {
  name: MODALS.statsDateInterval,
  components: {
    DatePicker,
    DateTimePicker,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      menu: false,
      form: {
        durationUnit: 'h',
        tstart: moment().startOf('hour'),
        tstopMode: false,
        tstop: '',
      },
      periodUnits: [
        {
          text: this.$tc('common.times.hour'),
          value: 'h',
        },
        {
          text: this.$tc('common.times.day'),
          value: 'd',
        },
        {
          text: this.$tc('common.times.week'),
          value: 'w',
        },
        {
          text: this.$tc('common.times.month'),
          value: 'm',
        },
      ],
    };
  },
  methods: {
    save(date) {
      this.$refs.menu.save(date);
    },
  },
};
</script>
