<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-text-field(
          v-model="form._id",
          :label="$t('modals.createWebhook.fields.id')",
          :readonly="isDisabledIdField",
          :disabled="isDisabledIdField"
        )
          v-tooltip(slot="append", left)
            v-icon(slot="activator") help_outline
            span {{ $t('modals.createWebhook.tooltips.id') }}
        retry-field(v-model="form.retry")
        v-switch(
          v-model="form.disable_if_active_pbehavior",
          :label="$t('webhook.disableIfActivePbehavior')"
        )
        webhook-form(v-model="form")
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

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import WebhookForm from '@/components/other/webhook/form/webhook-form.vue';
import RetryField from '@/components/forms/fields/retry.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: { WebhookForm, ModalWrapper, RetryField },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { webhook } = this.modal.config;
    const defaultForm = {
      hook: {
        triggers: [],
        event_patterns: [],
        alarm_patterns: [],
        entity_patterns: [],
      },
      request: {
        method: '',
        url: '',
        headers: [],
        payload: '{}',
      },
      declare_ticket: [],
      disable_if_active_pbehavior: false,
      emptyResponse: false,
    };

    return {
      form: webhook ? webhookToForm(webhook) : defaultForm,
    };
  },
  computed: {
    title() {
      if (this.config.webhook) {
        return this.$t('modals.createWebhook.edit.title');
      }

      return this.$t('modals.createWebhook.create.title');
    },
    isDisabledIdField() {
      return !!this.config.webhook;
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

          await this.config.action(formToWebhook(preparedForm));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
