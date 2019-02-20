<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Title
    v-card-text
      v-tabs(fixed-tabs)
        v-tab Hook
        v-tab Request
        v-tab(:disabled="hasBlockedTriggers") Declare ticket
        v-tab-item
          webhook-hook-tab(v-model="form.hook", :hasBlockedTriggers="hasBlockedTriggers")
        v-tab-item
          webhook-request-tab(v-model="form.request")
        v-tab-item
          webhook-declare-ticket-tab(v-model="form.declare_ticket")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, intersection } from 'lodash';

import { MODALS, WEBHOOK_TRIGGERS } from '@/constants';

import { setInSeveral } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';

import modalInnerMixin from '@/mixins/modal/inner';

import WebhookHookTab from './partials/webhook-hook-tab.vue';
import WebhookRequestTab from './partials/webhook-request-tab.vue';
import WebhookDeclareTicketTab from './partials/webhook-declare-ticket-tab.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWebhook,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    WebhookHookTab,
    WebhookRequestTab,
    WebhookDeclareTicketTab,
  },
  filters: {
    webhookToForm(webhook) {
      return setInSeveral(webhook, {
        'request.headers': objectToTextPairs,
      });
    },
    formToWebhook(form) {
      return setInSeveral(form, {
        'request.headers': textPairsToObject,
      });
    },
  },
  mixins: [modalInnerMixin],
  data() {
    const { webhook } = this.modal.config;
    const defaultForm = {
      hook: {
        triggers: [],
        event_pattern: [],
        alarm_pattern: [],
        entity_pattern: [],
      },
      request: {
        method: '',
        url: '',
        auth: {
          username: '',
          password: '',
        },
        headers: [],
        payload: '{}',
      },
      declare_ticket: {},
    };

    return {
      form: webhook ? this.$options.filters.webhookToForm(cloneDeep(webhook)) : defaultForm,
    };
  },
  computed: {
    hasBlockedTriggers() {
      return intersection(this.form.hook.triggers, [
        WEBHOOK_TRIGGERS.resolve,
        WEBHOOK_TRIGGERS.cancel,
        WEBHOOK_TRIGGERS.unsnooze,
      ]).length > 0;
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          const preparedForm = this.hasBlockedTriggers ? setInSeveral(this.form, {
            'hook.event_pattern': [],
            declare_ticket: {},
          }) : this.form;

          await this.config.action(this.$options.filters.formToWebhook(preparedForm));
        }

        this.hideModal();
      }
    },
  },
};
</script>
