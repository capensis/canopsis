<template lang="pug">
  transition(name="fade")
    div(v-show="pending")
      slot
</template>

<script>
export default {
  name: 'loader',
  props: {
    pending: {
      type: Boolean,
      required: true,
    },
  },
};
</script>
