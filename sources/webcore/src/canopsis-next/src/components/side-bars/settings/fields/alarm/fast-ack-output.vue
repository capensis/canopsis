<template lang="pug">
  v-list-group(data-test="fastAckOutput")
    v-list-tile(slot="activator") {{ $t('settings.fastAckOutput') }}
      .font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      enabled-field(v-field="value.enabled", data-test="fastAckOutputSwitch")
      v-text-field(
        v-field="value.value",
        data-test="fastAckOutputField",
        :label="$t('tables.alarmGeneral.output')",
        :disabled="!value.enabled"
      )
</template>

<script>
import EnabledField from '@/components/forms/fields/enabled-field.vue';

/**
 * Component to set fast-ack output
 *
 * @prop {Object} [value] - Value
 */
export default {
  components: { EnabledField },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
