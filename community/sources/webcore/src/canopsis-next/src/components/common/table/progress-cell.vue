<template>
  <td>
    <slot v-if="!pending" />
    <v-progress-circular
      class="ml-2"
      v-else
      indeterminate
      color="primary"
      width="3"
      size="25"
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
