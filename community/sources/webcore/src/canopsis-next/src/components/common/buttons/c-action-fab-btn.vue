<template lang="pug">
  c-action-btn(v-bind="$attrs", custom-tooltip-activator)
    template(#button="")
      v-btn(:color="color", :small="small", fab, dark, @click.stop="$listeners.click")
        v-icon {{ icon }}
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    icon: {
      type: String,
      required: true,
    },
    color: {
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
