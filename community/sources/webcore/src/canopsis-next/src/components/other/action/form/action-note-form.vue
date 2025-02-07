<template>
  <v-layout column>
    <c-payload-textarea-field
      v-field="value.output"
      :label="$t('scenario.output')"
      :variables="payloadVariables"
    >
      <template #append="">
        <c-help-icon
          :text="$t('scenario.outputHelp')"
          icon="help"
          left
        />
      </template>
    </c-payload-textarea-field>
  </v-layout>
</template>

<script>
import { payloadVariablesMixin } from '@/mixins/payload/variables';

export default {
  inject: ['$validator'],
  mixins: [payloadVariablesMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  computed: {
    payloadVariables() {
      return [
        ...this.alarmPayloadVariables,
        ...this.additionalDataVariables,
      ];
    },
  },
};
</script>
