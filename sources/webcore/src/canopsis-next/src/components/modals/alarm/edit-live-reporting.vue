<template lang="pug">
  v-card(data-test="liveReportingModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.liveReporting.editLiveReporting') }}
    v-card-text
      h3 {{ $t('modals.liveReporting.dateInterval') }}
      date-interval-selector(v-model="form")
      v-divider
      v-layout.py-1(justify-end)
        v-btn(
          @click="$modals.hide",
          depressed,
          flat,
          data-test="liveReportingCancelButton"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          @click="submit",
          :disabled="errors.any()",
          data-test="liveReportingApplyButton"
        ) {{ $t('common.apply') }}
</template>

<script>
import moment from 'moment';

import { MODALS, DATETIME_FORMATS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

/**
 * Modal to add a time filter on alarm-list
 */
export default {
  name: MODALS.editLiveReporting,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    DateIntervalSelector,
  },
  mixins: [modalInnerMixin],
  data() {
    const { config } = this.modal;

    return {
      form: {
        tstart: config.tstart || '',
        tstop: config.tstop || '',
      },
    };
  },
  computed: {
    tstartRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    tstopRules() {
      const rules = { required: true };

      if (this.tstart) {
        rules.after = [moment(this.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    },
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
