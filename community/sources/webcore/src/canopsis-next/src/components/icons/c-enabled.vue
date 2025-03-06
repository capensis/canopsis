<template>
  <span>
    <v-icon
      v-if="value"
      color="success"
    >
      {{ iconSuccess }}
    </v-icon>
    <v-icon
      v-else
      color="error"
    >
      {{ iconFailed }}
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
    iconSuccess: {
      type: String,
      default: 'check_circle',
    },
    iconFailed: {
      type: String,
      default: 'cancel',
    },
  },
};
</script>
