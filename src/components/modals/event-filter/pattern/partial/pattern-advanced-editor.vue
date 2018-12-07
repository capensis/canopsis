<template lang="pug">
  div
    v-textarea(:value="JSON.stringify(value, undefined, 4)", @input="checkValidity", rows="20")
    v-btn(:color="error ? 'error' : 'primary'") {{ error ? 'JSON not valid' : 'Valid JSON' }}
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
      error: '',
    };
  },
  methods: {
    checkValidity(event) {
      try {
        JSON.parse(event);
        this.$emit('input', JSON.parse(event));
        this.error = '';
      } catch (err) {
        this.error = err.message;
      }
    },
  },
};
</script>

