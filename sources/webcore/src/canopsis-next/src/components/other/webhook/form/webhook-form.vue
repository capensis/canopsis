<template lang="pug">
  v-tabs(fixed-tabs)
    v-tab {{ $t('webhook.tabs.hook.title') }}
    v-tab {{ $t('webhook.tabs.request.title') }}
    v-tab(:disabled="hasBlockedTriggers") {{ $t('webhook.tabs.declareTicket.title') }}
    v-tab-item
      webhook-form-hook-tab(
      :hook="form.hook",
      :hasBlockedTriggers="hasBlockedTriggers",
      :disabled="disabled",
      @input="updateField('hook', $event)"
      )
    v-tab-item
      webhook-form-request-tab(
      :request="form.request",
      :disabled="disabled",
      @input="updateField('request', $event)"
      )
    v-tab-item
      webhook-form-declare-ticket-tab(
      :declareTicket="form.declare_ticket",
      :disabled="disabled",
      @input="updateField('declare_ticket', $event)"
      )
</template>

<script>
import { intersection } from 'lodash';

import { WEBHOOK_TRIGGERS } from '@/constants';

import formMixin from '@/mixins/form';

import WebhookFormHookTab from './tabs/webhook-form-hook-tab.vue';
import WebhookFormRequestTab from './tabs/webhook-form-request-tab.vue';
import WebhookFormDeclareTicketTab from './tabs/webhook-form-declare-ticket-tab.vue';

export default {
  components: {
    WebhookFormHookTab,
    WebhookFormRequestTab,
    WebhookFormDeclareTicketTab,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
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
};
</script>
