<template lang="pug">
  transition(name="fade")
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
