<template>
  <v-text-field
    v-validate="rules"
    v-field="value"
    v-bind="$attrs"
    :label="label || $t('common.id')"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    :readonly="disabled"
    :name="name"
    @input="errors.remove(name)"
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
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: '_id',
    },
    helpText: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
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
