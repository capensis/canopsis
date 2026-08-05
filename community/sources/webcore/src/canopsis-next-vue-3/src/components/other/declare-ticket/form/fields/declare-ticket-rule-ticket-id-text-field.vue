<template>
  <c-text-field
    :model-value="value"
    :label="$t('declareTicket.ticketID')"
    :name="name"
    :disabled="disabled"
    :required="required"
    class="declare-ticket-rule-ticket-id-text-field"
    @update:model-value="$emit('update:value', $event)"
  >
    <template #append-inner>
      <c-help-icon
        :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketID') })"
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
      default: 'ticket_id',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.declare-ticket-rule-ticket-id-text-field {
  pointer-events: auto !important;
}
</style>
