<template lang="pug">
  v-card(data-test="statsDateIntervalModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.statsDateInterval.title') }}
    v-card-text
      v-container
        v-layout
          v-flex(xs3, v-if="!hiddenFields.includes('periodValue')")
            v-text-field.pt-0(
              type="number",
              v-model="periodForm.periodValue",
              :label="$t('modals.statsDateInterval.fields.periodValue')"
            )
          v-flex
            v-select.pt-0(
              v-model="periodForm.periodUnit",
              :items="periodUnits",
              :label="$t('modals.statsDateInterval.fields.periodUnit')"
            )
        v-alert.mb-2(
          :value="isPeriodMonth",
          type="info"
        ) {{ $t('settings.statsDateInterval.monthPeriodInfo') }}
        date-interval-selector.my-1(
          v-model="dateSelectorForm",
          tstopRules="after_custom:tstart",
          @update:startObjectValue="updateStartObjectValue",
          @update:stopObjectValue="updateStopObjectValue"
        )
        v-alert.mb-2(
          :value="isPeriodMonth",
          type="info"
        ) {{ monthIntervalMessage }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit", :disabled="errors.any()") {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS, DATETIME_FORMATS, STATS_DURATION_UNITS, STATS_QUICK_RANGES } from '@/constants';

import {
  dateParse,
  prepareStatsStopForMonthPeriod,
  prepareStatsStartForMonthPeriod,
} from '@/helpers/date-intervals';

import modalInnerMixin from '@/mixins/modal/inner';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

export default {
  name: MODALS.statsDateInterval,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DateIntervalSelector,
  },
  mixins: [modalInnerMixin],
  data() {
    const { interval } = this.modal.config;
    const defaultPeriodForm = {
      periodValue: 1,
      periodUnit: STATS_DURATION_UNITS.hour,
    };

    const defaultDateSelectorForm = {
      tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
      tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
    };

    let periodForm;
    let dateSelectorForm;

    if (interval) {
      periodForm = pick(interval, Object.keys(defaultPeriodForm));
      dateSelectorForm = pick(interval, Object.keys(defaultDateSelectorForm));
    } else {
      periodForm = defaultPeriodForm;
      dateSelectorForm = defaultDateSelectorForm;
    }

    return {
      periodForm,
      dateSelectorForm,

      dateObjectValues: {
        start: null,
        stop: null,
      },
    };
  },
  computed: {
    periodUnits() {
      return [
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
      ];
    },

    hiddenFields() {
      return this.modal.config.hiddenFields || [];
    },

    isPeriodMonth() {
      return this.periodForm.periodUnit === STATS_DURATION_UNITS.month;
    },

    monthIntervalMessage() {
      return this.$t('modals.statsDateInterval.info.monthPeriodUnit', {
        start: this.$options.filters.date(this.dateObjectValues.start, 'long', false),
        stop: this.$options.filters.date(this.dateObjectValues.stop, 'long', false),
      });
    },
  },
  created() {
    this.$validator.extend('after_custom', {
      getMessage: () => this.$t('modals.statsDateInterval.errors.endDateLessOrEqualStartDate'),
      validate: (value, [otherValue]) => {
        try {
          const convertedStop = dateParse(value, 'stop', DATETIME_FORMATS.dateTimePicker);
          const convertedStart = dateParse(otherValue, 'start', DATETIME_FORMATS.dateTimePicker);

          return !convertedStop.isSameOrBefore(convertedStart);
        } catch (err) {
          return true; // TODO: problem with i18n: https://github.com/baianat/vee-validate/issues/2025
        }
      },
    }, {
      hasTarget: true,
    });
  },
  methods: {
    updateStartObjectValue(value) {
      this.dateObjectValues.start = value && prepareStatsStartForMonthPeriod(value);
    },

    updateStopObjectValue(value) {
      this.dateObjectValues.stop = value && prepareStatsStopForMonthPeriod(value);
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          this.config.action({ ...this.periodForm, ...this.dateSelectorForm });
        }

        this.hideModal();
      }
    },
  },
};
</script>
