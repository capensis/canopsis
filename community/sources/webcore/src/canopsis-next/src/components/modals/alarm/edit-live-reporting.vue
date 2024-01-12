<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.liveReporting.editLiveReporting') }}</span>
      </template>
      <template #text="">
        <h3>{{ $t('modals.liveReporting.dateInterval') }}</h3>
        <date-interval-selector
          v-model="form"
          :quick-ranges="quickRanges"
        />
      </template>
      <template #actions="">
        <v-btn
          text
          depressed
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :loading="submitting"
          :disabled="isDisabled"
          type="submit"
        >
          {{ $t('common.apply') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, LIVE_REPORTING_QUICK_RANGES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

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
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { config } = this.modal;

    return {
      form: {
        tstart: config.tstart ?? '',
        tstop: config.tstop ?? '',
        time_field: config.time_field ?? '',
      },
    };
  },
  computed: {
    quickRanges() {
      return Object.values(LIVE_REPORTING_QUICK_RANGES);
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
