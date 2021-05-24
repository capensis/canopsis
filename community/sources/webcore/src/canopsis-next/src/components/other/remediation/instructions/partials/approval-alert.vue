<template lang="pug">
  v-card
    v-card-text.approval-alert-text
      v-icon.notification-icon(color="error") notification_important
      v-layout(row)
        strong {{ approval.requested_by }}
        span.ml-1 {{ $t('modals.remediationInstructionApproval.requested') }}
      v-layout.mt-3(row)
        span.grey--text.text--darken-2 {{ $tc('common.comment', 1) }}:
      v-layout.mt-3(row)
        span {{ approval.comment }}
</template>

<script>
export default {
  props: {
    approval: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.approval-alert-text {
  position: relative;
  padding-left: 60px;

  .notification-icon {
    position: absolute;
    left: 20px;
    top: 16px;
  }
}
</style>
