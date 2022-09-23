<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.fastAckOutput') }}
        div.font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      c-enabled-field(v-field="value.enabled")
      v-text-field(
        v-field="value.value",
        :label="$t('common.output')",
        :disabled="!value.enabled"
      )
</template>

<script>
/**
 * Component to set fast-ack output
 *
 * @prop {Object} [value] - Value
 */
export default {
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
