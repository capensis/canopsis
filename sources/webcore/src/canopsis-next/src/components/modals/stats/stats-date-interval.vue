<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Stats - Date interval
    v-card-text
      v-container
        v-layout
          v-flex(xs3)
            v-text-field.pt-0(
            type="number",
            v-model="periodValue",
            label="Period value"
            )
          v-select.pt-0(
          :items="periodUnits",
          v-model="periodUnit",
          label="Period unit",
          )
        v-alert.mb-2(
        v-if="periodUnit === 'm'", type="info", value="true"
        ) {{ $t('settings.statsDateInterval.monthPeriodInfo') }}
        stats-date-selector.my-1(v-model="dateForm", :periodUnit="periodUnit", @input="resetValidation")
      v-alert(
      value="errors",
      type="error",
      v-for="error in errors",
      :key="error",
      ) {{ error }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="errors.length !== 0") {{ $t('common.submit') }}
</template>

<script>
import moment from 'moment';

import { MODALS } from '@/constants';

import { parseStringToDateInterval } from '@/helpers/date-intervals';

import modalInnerMixin from '@/mixins/modal/inner';

import StatsDateSelector from '@/components/forms/stats-date-selector.vue';

export default {
  name: MODALS.statsDateInterval,
  components: {
    StatsDateSelector,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      periodValue: 1,
      periodUnit: 'h',
      dateForm: {
        tstart: 'now+1d',
        tstop: 'now+2d',
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
      errors: [],
    };
  },
  mounted() {
    if (this.config.interval) {
      const {
        periodValue,
        periodUnit,
        tstart,
        tstop,
      } = this.config.interval;

      this.periodValue = periodValue;
      this.periodUnit = periodUnit;
      this.dateForm = { ...this.dateForm, tstart, tstop };
    }
  },
  methods: {
    resetValidation() {
      this.errors = [];
    },

    validate() {
      const { tstart, tstop } = this.dateForm;

      try {
        const convertedTstart = moment(tstart).isValid() ? moment(tstart) : parseStringToDateInterval(tstart, 'start');
        const convertedTstop = moment(tstop).isValid() ? moment(tstop) : parseStringToDateInterval(tstop, 'stop');

        if (convertedTstop.isBefore(convertedTstart)) {
          this.errors.push('Tstop before tstart');
          return false;
        }
      } catch (err) {
        this.errors.push(err.message);
        return false;
      }

      return true;
    },

    async submit() {
      if (this.validate()) {
        if (this.config.action) {
          this.config.action({
            periodValue: this.periodValue,
            periodUnit: this.periodUnit,
            tstart: this.dateForm.tstart,
            tstop: this.dateForm.tstop,
          });
        }

        this.hideModal();
      }
    },
  },
};
</script>
