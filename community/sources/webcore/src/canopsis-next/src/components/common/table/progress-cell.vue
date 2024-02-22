<template>
  <td>
    <slot v-if="!pending" />
    <v-progress-circular
      v-else
      class="ml-2"
      color="primary"
      width="3"
      size="25"
      indeterminate
    />
  </td>
</template>

<script>
export default {
  props: {
    pending: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
