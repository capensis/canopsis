<template lang="pug">
  transition
    div(v-show="pending")
      slot
</template>

<script>
export default {
  props: {
    pending: {
      type: Boolean,
      required: true,
    },
  },
};
</script>
