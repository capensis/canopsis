<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(text-class="pa-0", close)
      template(#title="")
        span {{ title }}
      template(#text="")
        v-layout(column)
          broadcast-message(:message="message", :color="form.color")
          broadcast-message-form.pa-3(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary.white--text(
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { messageToForm, formToMessage } from '@/helpers/forms/broadcast-message';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';
import BroadcastMessageForm from '@/components/other/broadcast-message/broadcast-message-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createBroadcastMessage,
  $_veeValidate: {
    validator: 'new',
  },
  components: { BroadcastMessage, BroadcastMessageForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: messageToForm(this.modal.config.message),
    };
  },
  computed: {
    title() {
      return this.modal.config.title || this.$t('modals.createBroadcastMessage.create.title');
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
