<template>
  <v-text-field
    v-field="value"
    v-validate="rules"
    v-bind="$attrs"
    :label="label || $t('common.name')"
    :error-messages="errors.collect(name)"
    :name="name"
  />
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
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'name',
    },
    required: {
      type: Boolean,
      default: false,
    },
    maxLength: {
      type: Number,
      required: false,
    },
  },
  computed: {
    rules() {
      const rules = {
        required: this.required,
      };

      if (this.maxLength) {
        rules.max = this.maxLength;
      }

      return rules;
    },
  },
};
</script>
