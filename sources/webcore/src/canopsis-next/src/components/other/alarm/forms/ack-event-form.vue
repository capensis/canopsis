<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.ticket",
        :label="$t('modals.createAckEvent.fields.ticket')",
        data-test="createAckEventTicket"
      )
    v-layout(row)
      v-textarea(
        v-field="form.output",
        v-validate="isNoteRequired ? 'required' : ''",
        :label="$t('modals.createAckEvent.fields.output')",
        :error-messages="errors.collect('output')",
        name="output",
        data-test="createAckEventNote"
      )
    v-layout(row)
      v-tooltip(top)
        v-checkbox(
          slot="activator",
          v-field="form.ack_resources",
          :label="$t('modals.createAckEvent.fields.ackResources')",
          color="primary",
          data-test="createAckEventResource"
        )
        span {{ $t('modals.createAckEvent.tooltips.ackResources') }}
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNoteRequired: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
