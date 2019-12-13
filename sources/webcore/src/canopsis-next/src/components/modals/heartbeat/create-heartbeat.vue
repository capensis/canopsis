<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createHeartbeat.create.title') }}
      template(slot="text")
        heartbeat-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import HeartbeatForm from '@/components/other/heartbeat/form/heartbeat-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createHeartbeat,
  $_veeValidate: {
    validator: 'new',
  },
  components: { HeartbeatForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    return {
      form: {
        pattern: {},
        periodValue: '',
        periodUnit: '',
      },
    };
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const { pattern, periodValue, periodUnit } = this.form;
        const data = {
          pattern,
          expected_interval: `${periodValue}${periodUnit}`,
        };

        if (this.config.action) {
          await this.config.action(data);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
