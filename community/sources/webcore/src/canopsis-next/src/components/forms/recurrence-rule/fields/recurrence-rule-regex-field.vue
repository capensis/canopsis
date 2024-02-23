<template>
  <v-text-field
    v-field="value"
    v-validate="{ regex: advancedFieldRegex }"
    v-bind="$attrs"
    :error-messages="errors.collect(name)"
    :name="name"
    persistent-hint
  >
    <template #append="">
      <c-help-icon
        v-if="helpText"
        :text="helpText"
        icon="help"
        max-width="250"
        left
      />
    </template>
  </v-text-field>
</template>

<script>
export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'regex',
    },
  },
  computed: {
    advancedFieldRegex() {
      return /^(-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))(,-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))*,?$/;
    },
  },
};
</script>
