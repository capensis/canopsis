<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        template(v-if="enabled")
          v-layout(justify-center, align-center)
            v-progress-circular(:value="progress", :indeterminate="cancelling", color="primary")
            span.ml-4(v-if="!cancelling") {{ $t('maintenance.timerMessage', { timer }) }}
        template(v-else)
          c-alert(type="warning") {{ $t('maintenance.logoutWarning')}}
          maintenance-form(v-model="form")
      template(#actions="")
        template(v-if="!enabled")
          v-btn(
            depressed,
            flat,
            @click="$modals.hide"
          ) {{ $t('common.cancel') }}
          v-btn.primary(
            :loading="submitting",
            :disabled="isDisabled",
            type="submit"
          ) {{ $t('modals.createMaintenance.enableMaintenance') }}
        v-btn(
          v-else,
          :loading="cancelling",
          :disabled="isDisabledCancel",
          depressed,
          flat,
          @click="cancel"
        ) {{ $t('modals.createMaintenance.cancelMaintenance') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { maintenanceToForm } from '@/helpers/entities/maintenance/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import MaintenanceForm from '@/components/other/maintenance/form/maintenance-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMaintenance,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { MaintenanceForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    submittableMixinCreator({
      method: 'cancel',
      property: 'cancelling',
      computedProperty: 'isDisabledCancel',
    }),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { maintenance, cancelTimer } = this.modal.config;

    return {
      form: maintenanceToForm(maintenance),
      enabled: false,
      timer: cancelTimer,
      intervalId: null,
    };
  },
  computed: {
    title() {
      return this.$t(`modals.createMaintenance.${this.enabled ? 'cancel' : 'setup'}.title`);
    },

    progress() {
      return 100 - ((this.timer / this.config.cancelTimer) * 100);
    },
  },
  methods: {
    cancelTimer() {
      clearInterval(this.intervalId);
    },

    async logout() {
      await this.config.logout?.();
      this.$modals.hide();
    },

    startTimer() {
      this.intervalId = setInterval(() => {
        if (!this.timer) {
          this.cancelTimer();
          this.logout();
          return;
        }

        this.timer -= 1;
      }, 1000);
    },

    async cancel() {
      this.cancelTimer();

      await this.config.cancel?.();

      this.$modals.hide();
    },

    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        await this.config.action?.(this.form);

        if (this.config.cancelTimer > 0) {
          this.enabled = true;

          this.startTimer();
        } else {
          this.$modals.hide();
        }
      }
    },
  },
};
</script>
