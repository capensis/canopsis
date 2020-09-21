<template lang="pug">
  v-form(data-test="statsDateIntervalModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.statsDateInterval.title') }}
      template(slot="text")
        stats-date-interval-form(v-model="form", :hiddenFields="config.hiddenFields")
      template(slot="actions")
        v-btn(
          data-test="statsDateIntervalCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="statsDateIntervalSubmitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS, STATS_DURATION_UNITS, STATS_QUICK_RANGES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import StatsDateIntervalForm from '@/components/other/stats/stats-date-interval-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.statsDateInterval,
  $_veeValidate: {
    validator: 'new',
  },
  components: { StatsDateIntervalForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
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
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
