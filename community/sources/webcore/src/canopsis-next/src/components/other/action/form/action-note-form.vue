<template>
  <v-layout column>
    <c-payload-textarea-field
      v-field="value.output"
      :label="$t('scenario.output')"
      :variables="templateVars.output"
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
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
