<template>
  <c-text-field
    :model-value="value"
    :label="$t('declareTicket.ticketURLTitle')"
    @update:model-value="$emit('update:value', $event)"
  >
    <template #append-inner>
      <c-help-icon
        :text="$t('declareTicket.ticketURLTitleHelp')"
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
  },
};
</script>
