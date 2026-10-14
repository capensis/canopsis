<template>
  <c-text-field
    :model-value="value"
    :label="$t('declareTicket.ticketID')"
    :name="name"
    :disabled="disabled"
    :required="required"
    :validation-rules="validationRules"
    class="declare-ticket-rule-ticket-id-text-field"
    @update:model-value="$emit('update:value', $event)"
  >
    <template #append-inner>
      <c-help-icon
        icon="help"
        color="grey darken-1"
        left
      >
        <div>
          {{ $t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketID') }) }}
        </div>
        <template v-if="ticketFormats.length">
          <div class="mt-2">{{ $t('declareTicket.allowedTicketFormats') }}</div>
          <code
            v-for="(format, index) in ticketFormats"
            :key="`${format}-${index}`"
            class="d-block"
          >{{ format }}</code>
        </template>
      </c-help-icon>
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
    validationRules: {
      type: Array,
      default: () => [],
    },
    ticketFormats: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

<style lang="scss">
.declare-ticket-rule-ticket-id-text-field {
  pointer-events: auto !important;
}
</style>
