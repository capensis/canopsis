<template lang="pug">
  v-text-field(
    v-field="value",
    v-validate="rules",
    :label="label || $t('common.name')",
    :error-messages="errors.collect(name)",
    :name="name",
    :disabled="disabled",
    :browser-autocomplete="browserAutocomplete"
  )
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
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    browserAutocomplete: {
      type: String,
      required: false,
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
