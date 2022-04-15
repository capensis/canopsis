<template lang="pug">
  c-action-btn(:tooltip="tooltip", top)
    template(#button="")
      v-btn(:color="color", :small="small", fab, dark, @click.stop="$listeners.click")
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
    small: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
