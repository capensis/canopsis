<template>
  <c-card-iterator-form
    :value="webhooks"
    :disabled="disabled"
    :draggable-group="draggableGroup"
    :name="name"
    :required-error-message="$t('declareTicket.errors.webhookRequired')"
    :empty-message="$t('declareTicket.emptyWebhooks')"
    :add-button-label="$t('declareTicket.addWebhook')"
    item-key="key"
    required
    @update:value="updateWebhooks"
    @add="addWebhook"
  >
    <template #item="{ index, item: webhook }">
      <declare-ticket-rule-webhook-field
        :form="webhooks[index]"
        :name="`${name}.${webhook.key}`"
        :is-declare-ticket-exist="!webhook.declare_ticket.enabled && isSomeOneDeclareTicketEnabled"
        :has-previous="!!index"
        :webhook-number="index + 1"
        :template-vars="templateVars"
        @update:form="updateItem(index, $event)"
        @remove="removeItemFromArray(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { computed } from 'vue';

import { declareTicketRuleWebhookToForm } from '@/helpers/entities/declare-ticket/rule/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import DeclareTicketRuleWebhookField from './declare-ticket-rule-webhook-field.vue';

export default {
  components: {
    DeclareTicketRuleWebhookField,
  },

  props: {
    webhooks: {
      type: Array,
      default: () => [],
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

  emits: ['update:webhooks'],

  setup(props, { emit }) {
    const { updateModel, addItemIntoArray, removeItemFromArray, updateItemInArray } =
      useArrayModelField(props, emit, 'webhooks');

    const isSomeOneDeclareTicketEnabled = computed(() =>
      props.webhooks.some(webhook => webhook.declare_ticket.enabled),
    );

    const draggableGroup = computed(() => ({
      name: 'declare-ticket-steps',
    }));

    const addWebhook = () => addItemIntoArray(declareTicketRuleWebhookToForm());
    const updateWebhooks = value => updateModel(value);
    const updateItem = (index, value) => updateItemInArray(index, value);

    return {
      isSomeOneDeclareTicketEnabled,
      draggableGroup,
      addWebhook,
      updateWebhooks,
      updateItem,
      removeItemFromArray,
    };
  },
};
</script>
