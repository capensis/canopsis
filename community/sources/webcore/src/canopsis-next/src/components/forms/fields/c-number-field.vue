<template>
  <v-text-field
    v-bind="$attrs"
    v-field.number="value"
    v-validate="rules"
    :label="label"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    :hide-details="hideDetails"
    :name="name"
    :min="min"
    :max="max"
    :step="step"
    :loading="loading"
    type="number"
    @paste="$emit('paste', $event)"
    @click="$emit('click', $event)"
  >
    <template #append="">
      <slot name="append" />
    </template>
  </v-text-field>
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
      type: [Number, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'value',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: undefined,
    },
    max: {
      type: Number,
      default: undefined,
    },
    step: {
      type: Number,
      default: undefined,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
        numeric: !this.step,
        decimal: !!this.step,
        min_value: this.min ?? false,
        max_value: this.max ?? false,
      };
    },
  },
};
</script>
