<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
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

import { heartbeatToForm, formToHeartbeat } from '@/helpers/forms/heartbeat';

import authMixin from '@/mixins/auth';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

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
  mixins: [
    authMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { heartbeat = {} } = this.modal.config;

    return {
      form: heartbeatToForm(heartbeat),
    };
  },
  computed: {
    title() {
      let type = 'create';

      if (this.config.heartbeat) {
        type = this.config.isDuplicating ? 'duplicate' : 'edit';
      }

      return this.$t(`modals.createHeartbeat.${type}.title`);
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          const heartbeat = formToHeartbeat(this.form);

          heartbeat.author = this.currentUser._id;

          await this.config.action(heartbeat);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
