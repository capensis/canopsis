<template lang="pug">
  component(
    v-clipboard:copy="value",
    v-clipboard:success="showSuccessPopup",
    v-clipboard:error="showErrorPopup",
    :is="tag"
  )
    slot
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  computed: {
    tag() {
      const [slot] = this.$slots.default || [];

      return slot?.tag ?? 'span';
    },
  },
  methods: {
    showSuccessPopup() {
      this.$popups.success({ text: this.$t('popups.copySuccess') });
    },

    showErrorPopup() {
      this.$popups.error({ text: this.$t('popups.copyError') });
    },
  },
};
</script>
