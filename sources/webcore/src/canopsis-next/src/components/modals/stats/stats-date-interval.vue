<template lang="pug">
  v-card(data-test="statsDateIntervalModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.statsDateInterval.title') }}
    v-card-text
      stats-date-interval-form(v-model="form", :hiddenFields="config.hiddenFields")
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(
          data-test="statsDateIntervalCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          data-test="statsDateIntervalSubmitButton",
          :disabled="errors.any()",
          @click="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS, STATS_DURATION_UNITS, STATS_QUICK_RANGES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import StatsDateIntervalForm from '@/components/other/stats/stats-date-interval-form.vue';

export default {
  name: MODALS.statsDateInterval,
  $_veeValidate: {
    validator: 'new',
  },
  components: { StatsDateIntervalForm },
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
      form: {
        ...periodForm,
        ...dateSelectorForm,
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
