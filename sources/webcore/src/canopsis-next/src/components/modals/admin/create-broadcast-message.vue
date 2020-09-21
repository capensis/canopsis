<template lang="pug">
  v-form.create-broadcast-message-modal(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        broadcast-message(:message="message", :color="form.color")
        broadcast-message-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary.white--text(
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { messageToForm, formToMessage } from '@/helpers/forms/broadcast-message';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';
import BroadcastMessageForm from '@/components/other/broadcast-message/broadcast-message-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createBroadcastMessage,

  $_veeValidate: {
    validator: 'new',
  },
  components: { BroadcastMessage, BroadcastMessageForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { message = {} } = this.modal.config;

    return {
      form: messageToForm(message),
    };
  },
  computed: {
    title() {
      const type = this.modal.config.message ? 'edit' : 'create';

      return this.$t(`modals.createBroadcastMessage.${type}.title`);
    },
    message() {
      return this.form.message || this.$t('modals.createBroadcastMessage.defaultMessage');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToMessage(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .create-broadcast-message-modal {
    & /deep/ .v-card__text {
      position: relative;
    }

    & /deep/ .broadcast-message {
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
    }
  }
</style>
