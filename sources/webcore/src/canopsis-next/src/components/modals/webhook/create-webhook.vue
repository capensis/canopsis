<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      v-text-field(
        v-model="form._id",
        :label="$t('modals.createWebhook.fields.id')",
        :readonly="isDisabledIdField",
        :disabled="isDisabledIdField"
      )
        v-tooltip(slot="append", left)
          v-icon(slot="activator") help_outline
          span {{ $t('modals.createWebhook.tooltips.id') }}
      v-switch(
        v-model="form.disable_if_active_pbehavior",
        :label="$t('webhook.disableIfActivePbehavior')"
      )
      webhook-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';
import { formToWebhook, webhookToForm } from '@/helpers/forms/webhook';

import modalInnerMixin from '@/mixins/modal/inner';

import WebhookForm from '@/components/other/webhook/form/webhook-form.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: { WebhookForm },
  mixins: [modalInnerMixin],
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

        this.hideModal();
      }
    },
  },
};
</script>
