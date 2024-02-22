<template>
  <v-textarea
    v-field="value"
    v-validate="rules"
    :label="label || $t('common.description')"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    :name="name"
  >
    <template
      v-if="helpText"
      #append=""
    >
      <c-help-icon
        :text="helpText"
        icon="help"
        left
      />
    </template>
  </v-textarea>
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
      default: 'description',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    helpText: {
      type: String,
      default: '',
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
