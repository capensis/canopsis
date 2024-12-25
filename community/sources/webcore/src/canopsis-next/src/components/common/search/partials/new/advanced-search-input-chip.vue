<template>
  <v-chip
    :close="closable"
    @click:close="close"
  >

  </v-chip>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    prefix: {
      type: String,
      default: '',
    },
    closable: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const close = () => {};

    return {
      close,
    };
  },
};
</script>
