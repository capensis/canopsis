<template>
  <v-card>
    <v-card-text class="approval-alert-text">
      <v-icon
        class="notification-icon"
        color="error"
      >
        notification_important
      </v-icon>
      <v-layout>
        <strong>{{ approval.requested_by.name }}</strong>
        <span class="ml-1">
          {{ $t('modals.remediationInstructionApproval.requested') }}
        </span>
      </v-layout>
      <v-layout
        class="mt-3"
      >
        <span class="grey--text text--darken-2">{{ $tc('common.comment', 1) }}:</span>
      </v-layout>
      <v-layout
        class="mt-3"
      >
        <span>{{ approval.comment }}</span>
      </v-layout>
    </v-card-text>
  </v-card>
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
