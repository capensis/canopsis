<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        webhook-form(v-model="form", :is-disabled-id-field="isDisabledIdField")
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

import { setSeveralFields } from '@/helpers/immutable';
import { formToWebhook, webhookToForm } from '@/helpers/forms/webhook';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import WebhookForm from '@/components/other/webhook/form/webhook-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    WebhookForm,
    ModalWrapper,
  },
  mixins: [
    authMixin,
    modalInnerMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { webhook } = this.modal.config;

    return {
      form: webhookToForm(webhook),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createWebhook.create.title');
    },

    isDisabledIdField() {
      return this.config.isDisabledIdField;
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          const preparedForm = this.hasBlockedTriggers ? setSeveralFields(this.form, {
            'hook.event_patterns': null,
            declare_ticket: {},
          }) : this.form;

          const webhook = formToWebhook(preparedForm);
          webhook.author = this.currentUser._id;

          await this.config.action(webhook);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
