<template lang="pug">
  v-card
    v-card-text
      v-layout.approval-alert-text(row)
        div.approval-alert-text__icon
          v-icon(color="error") notification_important
        v-layout(column)
          v-layout(row)
            strong {{ userName }}
            span.ml-1 {{ statusText }}
          v-layout.mt-3(row)
            span.grey--text.text--darken-2 {{ $tc('common.comment', 1) }}:
          v-layout.mt-3(row)
            span {{ comment }}
          slot
</template>

<script>
export default {
  props: {
    userName: {
      type: String,
      required: true,
    },
    comment: {
      type: String,
      required: false,
    },
    dismissed: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    statusText() {
      return this.$t(`modals.remediationInstructionApproval.${this.dismissed ? 'dismissed' : 'requested'}`);
    },
  },
};
</script>

<style lang="scss" scoped>
.approval-alert-text {
  &__icon {
    display: flex;
    justify-content: center;
    align-items: flex-start;
    width: 60px;
    flex-shrink: 0 !important;
  }
}
</style>
