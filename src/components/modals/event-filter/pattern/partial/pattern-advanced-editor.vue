<template lang="pug">
  div
    v-textarea(:value="JSON.stringify(value, undefined, 4)", @input="checkValidity", rows="15")
    v-btn(
    :color="error ? 'error' : 'primary'",
    :disabled="error ? true : false",
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
      error: '',
    };
  },
  methods: {
    checkValidity(event) {
      try {
        this.newVal = JSON.parse(event);
        this.error = '';
      } catch (err) {
        this.error = err.message;
      }
    },
    save() {
      this.$emit('input', this.newVal);
    },
  },
};
</script>

