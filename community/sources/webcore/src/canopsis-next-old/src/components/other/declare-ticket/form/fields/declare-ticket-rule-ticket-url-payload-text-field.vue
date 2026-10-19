<template>
  <c-payload-text-field
    v-field="value"
    :label="$t('declareTicket.ticketURL')"
    :name="name"
    :variables="variables"
    :disabled="disabled"
  >
    <template #append="">
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
  model: {
    prop: 'value',
    event: 'input',
  },
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
