<template>
  <v-text-field
    class="declare-ticket-rule-ticket-id-field"
    v-field="value"
    v-validate="rules"
    :label="$t('declareTicket.ticketID')"
    :error-messages="errors.collect(name)"
    :name="name"
    :disabled="disabled"
  >
    <template #append="">
      <c-help-icon
        :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketID') })"
        icon="help"
        color="grey darken-1"
        left
      />
    </template>
  </v-text-field>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
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
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>

<style lang="scss">
.declare-ticket-rule-ticket-id-field {
  pointer-events: auto !important;
}
</style>
