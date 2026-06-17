<template>
  <c-information-block :title="$t('declareTicket.ticketStatusSourceField')">
    <v-layout class="gap-3">
      <v-flex xs6>
        <c-payload-text-field
          v-if="form.template"
          v-field="form.value"
          :label="$t('declareTicket.ticketStatusSourceField')"
          :name="name"
          :variables="variables"
          :disabled="disabled"
          required
        />
        <c-name-field
          v-else
          v-field="form.value"
          :label="$t('declareTicket.ticketStatusSourceField')"
          :name="name"
          :disabled="disabled"
          required
        />
      </v-flex>
      <c-enabled-field
        v-field="form.template"
        :label="$t('declareTicket.allowTicketStatusSourceTemplate')"
      />
    </v-layout>
  </c-information-block>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'ticket_status_source',
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
