<template lang="pug">
  v-card(flat)
    v-card-text
      v-layout(row, align-center)
        span {{ $t('eventFilter.event') }}:
        c-copy-btn(
          :value="eventString",
          :tooltip="$t('eventFilter.copyEventToClipboard')",
          left,
          small,
          icon-small,
          @success="onSuccessCopied",
          @error="onErrorCopied"
        )
      c-json-treeview(:json="eventString")
</template>

<script>
export default {
  props: {
    failure: {
      type: Object,
      required: true,
    },
  },
  computed: {
    eventString() {
      return JSON.stringify(this.failure.event);
    },
  },
  methods: {
    onSuccessCopied() {
      this.$popups.success({ text: this.$t('eventFilter.eventCopied') });
    },

    onErrorCopied() {
      this.$popups.success({ text: this.$t('errors.default') });
    },
  },
};
</script>
