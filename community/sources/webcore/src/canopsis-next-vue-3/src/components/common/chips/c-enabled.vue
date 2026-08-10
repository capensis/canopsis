<template>
  <span>
    <c-icon v-if="value" color="success">
      {{ iconSuccess }}
    </c-icon>
    <c-icon v-else color="error">
      {{ iconFailed }}
    </c-icon>
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
