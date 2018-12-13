<template lang="pug">
  div
    v-textarea(:value="patternValue", @input="checkValidity", rows="15")
    v-btn(
    :color="error ? 'error' : 'primary'",
    :disabled="error",
    @click="save",
    ) {{ error ? $t('errors.JSONNotValid') : $t('common.save') }}
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      newVal: {},
      error: false,
      valueChanged: false,
    };
  },
  computed: {
    patternValue() {
      try {
        return JSON.stringify(this.value, undefined, 4);
      } catch (err) {
        console.error(err);
        return '';
      }
    },
  },
  methods: {
    checkValidity(value) {
      this.valueChanged = true;

      try {
        this.newVal = JSON.parse(value);
        this.error = false;
      } catch (err) {
        this.error = true;
      }
    },
    save() {
      if (this.valueChanged) {
        this.$emit('input', this.newVal);
      } else {
        this.$emit('input', this.value);
      }
    },
  },
};
</script>

