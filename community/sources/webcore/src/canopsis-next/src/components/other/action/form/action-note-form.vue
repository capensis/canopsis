<template>
  <v-layout column>
    <v-textarea
      v-field="value.output"
      :label="$t('scenario.output')"
    >
      <template #append="">
        <c-help-icon
          :text="$t('scenario.outputHelp')"
          icon="help"
          left
        />
      </template>
    </v-textarea>
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
  },
};
</script>
