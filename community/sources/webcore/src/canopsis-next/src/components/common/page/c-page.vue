<template lang="pug">
  div
    c-page-header
    v-card.ma-4.mt-0
      slot
    c-fab-btn(
      :has-access="creatable",
      :dark="dark",
      @refresh="$emit('refresh')",
      @create="$emit('create')"
    )
      span {{ createTooltip }}
</template>

<script>
export default {
  props: {
    creatable: {
      type: Boolean,
      default: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    createTooltip: {
      type: String,
      required: false,
    },
  },
};
</script>
