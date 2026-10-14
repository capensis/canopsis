<template>
  <c-text-field
    :model-value="value"
    :label="$t('declareTicket.ticketURL')"
    :name="name"
    :disabled="disabled"
    @update:model-value="$emit('update:value', $event)"
  >
    <template #append-inner>
      <c-help-icon
        :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketURL') })"
        icon="help"
        color="grey darken-1"
        left
      />
    </template>
  </c-text-field>
</template>

<script>
export default {
  emits: ['update:value'],
  props: {
    value: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'ticket_url',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
