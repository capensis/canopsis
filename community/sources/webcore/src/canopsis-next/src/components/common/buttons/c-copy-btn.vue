<template lang="pug">
  c-action-btn(v-bind="$attrs", :tooltip="tooltip")
    template(#button="")
      v-btn.mx-1(
        v-clipboard:copy="value",
        v-clipboard:success="onSuccessCopied",
        v-clipboard:error="onErrorCopied",
        icon
      )
        v-icon(:color="color") {{ icon }}
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      default: 'content_copy',
    },
    color: {
      type: String,
      default: '',
    },
    tooltip: {
      type: String,
      default: '',
    },
    value: {
      type: String,
      required: true,
    },
  },
  methods: {
    onSuccessCopied() {
      this.$emit('success');
    },

    onErrorCopied() {
      this.$emit('error');
    },
  },
};
</script>
