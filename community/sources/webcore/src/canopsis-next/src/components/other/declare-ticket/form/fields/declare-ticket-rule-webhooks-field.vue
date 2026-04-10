<template>
  <c-card-iterator-form
    v-field="webhooks"
    :disabled="disabled"
    :draggable-group="draggableGroup"
    :name="name"
    :required-error-message="$t('declareTicket.errors.webhookRequired')"
    :empty-message="$t('declareTicket.emptyWebhooks')"
    :add-button-label="$t('declareTicket.addWebhook')"
    iterator-class="mb-2"
    item-key="key"
    required
    @add="addWebhook"
  >
    <template #item="{ index, item: webhook }">
      <declare-ticket-rule-webhook-field
        v-field="webhooks[index]"
        :name="`${name}.${webhook.key}`"
        :is-declare-ticket-exist="!webhook.declare_ticket.enabled && isSomeOneDeclareTicketEnabled"
        :has-previous="!!index"
        :webhook-number="index + 1"
        :template-vars="templateVars"
        @remove="removeItemFromArray(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { declareTicketRuleWebhookToForm } from '@/helpers/entities/declare-ticket/rule/form';

import { formArrayMixin } from '@/mixins/form';

import DeclareTicketRuleWebhookField from './declare-ticket-rule-webhook-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleWebhookField,
  },
  mixins: [formArrayMixin],
  model: {
    prop: 'webhooks',
    event: 'input',
  },
  props: {
    webhooks: {
      type: Array,
      default: () => ([]),
    },
    name: {
      type: String,
      default: 'webhooks',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isSomeOneDeclareTicketEnabled() {
      return this.webhooks.some(webhook => webhook.declare_ticket.enabled);
    },

    hasWebhooksErrors() {
      return this.errors.has(this.name);
    },

    draggableGroup() {
      return {
        name: 'declare-ticket-steps',
      };
    },
  },
  methods: {
    addWebhook() {
      this.addItemIntoArray(declareTicketRuleWebhookToForm());
    },
  },
};
</script>
