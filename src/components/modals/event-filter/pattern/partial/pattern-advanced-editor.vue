<template lang="pug">
  div
    v-textarea(:value="JSON.stringify(value, undefined, 4)", @input="checkValidity")
    v-btn(:color="error ? 'error' : 'primary'") {{ error ? 'Not valid' : 'Valid' }}
</template>

<script>

export default {
  props: {
    value: {
      type: String,
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

