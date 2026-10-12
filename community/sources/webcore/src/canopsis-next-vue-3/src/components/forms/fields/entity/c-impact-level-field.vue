<template>
  <c-number-field
    :value="currentValue"
    :label="label || $t('common.impactLevel')"
    :placeholder="placeholder"
    :disabled="disabled"
    :name="name"
    :max="max"
    :min="min"
    :required="required"
    @update:value="updateValue"
  />
</template>

<script>
export default {
  emits: ['update:value', 'update:modelValue'],

  props: {
    value: {
      type: [Number, String],
      required: false,
    },
    modelValue: {
      type: [Number, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'impact_level',
    },
    min: {
      type: Number,
      default: 1,
    },
    max: {
      type: Number,
      default: 10,
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    currentValue() {
      return this.modelValue ?? this.value;
    },

    placeholder() {
      return `${this.min}-${this.max}`;
    },
  },

  methods: {
    updateValue(value) {
      this.$emit('update:value', value);
      this.$emit('update:modelValue', value);
    },
  },
};
</script>
