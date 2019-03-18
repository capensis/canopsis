<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.statsDateInterval.title') }}
    v-card-text
      v-container
        v-layout
          v-flex(xs3)
            v-text-field.pt-0(
            type="number",
            v-model="periodValue",
            :label="$t('modals.statsDateInterval.fields.periodValue')"
            )
          v-select.pt-0(
          v-model="periodUnit",
          :items="periodUnits",
          :label="$t('modals.statsDateInterval.fields.periodUnit')"
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
import { MODALS, DATETIME_FORMATS, STATS_DURATION_UNITS } from '@/constants';

import { dateParse } from '@/helpers/date-intervals';

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
      periodUnit: STATS_DURATION_UNITS.hour,
      dateForm: {
        tstart: 'now+1d',
        tstop: 'now+2d',
      },
      periodUnits: [
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
      this.dateForm = {
        tstart,
        tstop,
      };
    }
  },
  methods: {
    resetValidation() {
      this.errors = [];
    },

    validate() {
      const { tstart, tstop } = this.dateForm;

      try {
        const convertedTstart = dateParse(tstart, 'start', DATETIME_FORMATS.dateTimePicker);
        const convertedTstop = dateParse(tstop, 'stop', DATETIME_FORMATS.dateTimePicker);

        if (convertedTstop.isSameOrBefore(convertedTstart)) {
          this.errors.push(this.$t('modals.statsDateInterval.errors.tstopShouldBeMoreThanTstart'));
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
