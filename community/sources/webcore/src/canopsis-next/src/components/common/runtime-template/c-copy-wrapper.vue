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
      this.$popups.success({ text: this.$t('testSuite.popups.systemMessageCopied') });
    },

    showErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>
