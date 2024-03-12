<template>
  <v-text-field
    v-field="value"
    v-validate="rules"
    v-bind="$attrs"
    :label="$t('common.password')"
    :error-messages="errors.collect(name)"
    :name="name"
    type="password"
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
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'password',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>
