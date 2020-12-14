<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.ticket",
        :label="$t('modals.createAckEvent.fields.ticket')",
        data-test="createAckEventTicket"
      )
    output-field(v-field="form.output", :required="isNoteRequired")
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
import OutputField from '@/components/forms/fields/output.vue';

export default {
  components: { OutputField },
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
