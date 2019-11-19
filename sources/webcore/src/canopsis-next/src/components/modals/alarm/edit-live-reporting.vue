<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(data-test="liveReportingModal")
      template(slot="title")
        span {{ $t('modals.liveReporting.editLiveReporting') }}
      template(slot="text")
        h3 {{ $t('modals.liveReporting.dateInterval') }}
        date-interval-selector(v-model="form")
      template(slot="actions")
        v-btn(
          data-test="liveReportingCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="errors.any() || submitting",
          type="submit",
          data-test="liveReportingApplyButton"
        ) {{ $t('common.apply') }}
</template>

<script>
import moment from 'moment';

import { MODALS, DATETIME_FORMATS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to add a time filter on alarm-list
 */
export default {
  name: MODALS.editLiveReporting,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DateIntervalSelector, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    const { config } = this.modal;

    return {
      submitting: false,
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
      try {
        this.submitting = true;

        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          if (this.config.action) {
            await this.config.action(this.form);
          }

          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
