<template lang="pug">
  v-list-group(data-test="fastAckOutput")
    v-list-tile(slot="activator") {{ $t('settings.fastAckOutput') }}
      .font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      v-switch(
      color="primary",
      :label="$t('common.enabled')",
      :input-value="value.enabled",
      @change="updateField('enabled', $event)"
      )
      v-text-field(
      :label="$t('tables.alarmGeneral.output')",
      :value="value.value",
      :disabled="!value.enabled",
      @input="updateField('value', $event)"
      )
</template>

<script>
import formMixin from '@/mixins/form';

/**
* Component to set fast-ack output
*
* @prop {Object} [value] - Value
*/
export default {
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
