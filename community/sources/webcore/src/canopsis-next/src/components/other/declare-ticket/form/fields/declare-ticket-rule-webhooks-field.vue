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
import { computed } from 'vue';

import { declareTicketRuleWebhookToForm } from '@/helpers/entities/declare-ticket/rule/form';

import { useValidator } from '@/hooks/validator/validator';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import DeclareTicketRuleWebhookField from './declare-ticket-rule-webhook-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DeclareTicketRuleWebhookField,
  },
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
  setup(props, { emit }) {
    const validator = useValidator();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const isSomeOneDeclareTicketEnabled = computed(() => props.webhooks.some(
      webhook => webhook.declare_ticket.enabled,
    ));

    const hasWebhooksErrors = computed(() => validator?.errors?.has(props.name) ?? false);

    const draggableGroup = computed(() => ({
      name: 'declare-ticket-steps',
    }));

    const addWebhook = () => addItemIntoArray(declareTicketRuleWebhookToForm());

    return {
      isSomeOneDeclareTicketEnabled,
      hasWebhooksErrors,
      draggableGroup,
      addWebhook,
      removeItemFromArray,
    };
  },
};
</script>
