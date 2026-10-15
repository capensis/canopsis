<template>
  <c-payload-text-field
    :value="value"
    :label="$t('declareTicket.ticketURL')"
    :name="name"
    :variables="variables"
    :disabled="disabled"
    @update:value="$emit('update:value', $event)"
  >
    <template #append>
      <c-help-icon
        :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketURL') })"
        icon="help"
        color="grey darken-1"
        left
      />
    </template>
  </c-payload-text-field>
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
    variables: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
