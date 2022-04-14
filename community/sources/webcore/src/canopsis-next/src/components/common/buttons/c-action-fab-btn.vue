<template lang="pug">
  c-action-btn(:tooltip="tooltip", top)
    template(#button="")
      v-btn(:color="color", fab, dark, small, @click.stop="$listeners.click")
        v-icon {{ icon }}
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    color: {
      type: String,
      default: '',
    },
    tooltip: {
      type: String,
      default: '',
    },
  },
};
</script>
