<template lang="pug">
  v-text-field(
    v-validate="'numeric|min_value:0|max_value:100'",
    :label="label",
    :error-messages="errors.collect(name)",
    :value="percentValue",
    :name="name",
    :max="max",
    :min="min",
    type="number",
    prefix="%",
    @input="handleChange"
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
      type: Number,
      default: 0,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      required: false,
    },
    max: {
      type: Number,
      default: 100,
    },
    min: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    percentValue() {
      return Math.round(this.value * 100);
    },
  },
  methods: {
    handleChange(value) {
      this.$emit('input', Number(value) / 100);
    },
  },
};
</script>

