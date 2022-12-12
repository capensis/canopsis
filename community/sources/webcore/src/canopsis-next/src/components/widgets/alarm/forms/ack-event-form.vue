<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.ticket",
      :label="$t('modals.createAckEvent.fields.ticket')"
    )
    c-description-field(
      v-field="form.output",
      :label="$t('common.note')",
      :required="isNoteRequired"
    )
    v-layout(row)
      v-tooltip(top)
        template(#activator="{ on }")
          v-checkbox(
            v-on="on",
            v-field="form.ack_resources",
            :label="$t('modals.createAckEvent.fields.ackResources')",
            color="primary"
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
