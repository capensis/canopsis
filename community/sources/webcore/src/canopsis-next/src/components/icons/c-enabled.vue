<template>
  <span>
    <v-icon
      v-if="value"
      color="primary"
    >
      check_circle
    </v-icon>
    <v-icon
      v-else
      color="error"
    >
      cancel
    </v-icon>
  </span>
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
