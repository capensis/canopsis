<template lang="pug">
  v-layout(column)
    v-tooltip(v-if="!hideAckResources", top)
      template(#activator="{ on }")
        v-checkbox(
          v-on="on",
          v-field="form.ack_resources",
          :label="$t('alarm.ackResources')",
          color="primary"
        )
      span {{ $t('alarm.ackResourcesQuestion') }}
    c-description-field(
      v-field="form.output",
      :label="$t('common.note')",
      :required="isNoteRequired"
    )
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
    hideAckResources: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
