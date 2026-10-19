<template>
  <c-action-btn
    v-bind="$attrs"
    :tooltip="tooltip"
    :small="small"
    :icon-small="iconSmall"
    :icon="icon"
    :color="color"
    @click="copy"
  />
</template>

<script>
export default {
  inheritAttrs: false,
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
    small: {
      type: Boolean,
      default: false,
    },
    iconSmall: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['success', 'error'],
  methods: {
    async copy() {
      try {
        await navigator.clipboard.writeText(this.value);
        this.$emit('success');
      } catch (err) {
        this.$emit('error', err);
      }
    },
  },
};
</script>
