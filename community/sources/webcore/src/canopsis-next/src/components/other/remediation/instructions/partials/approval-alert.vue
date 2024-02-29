<template>
  <v-card>
    <v-card-text>
      <v-layout class="approval-alert-text">
        <div class="approval-alert-text__icon">
          <v-icon color="error">
            notification_important
          </v-icon>
        </div>
        <v-layout column>
          <v-layout>
            <strong>{{ userName }}</strong>
            <span class="ml-1">{{ statusText }}</span>
          </v-layout>
          <v-layout class="mt-3">
            <span class="grey--text text--darken-2">{{ $tc('common.comment', 1) }}:</span>
          </v-layout>
          <v-layout class="mt-3">
            <span>{{ comment }}</span>
          </v-layout>

          <slot />
        </v-layout>
      </v-layout>
    </v-card-text>
  </v-card>
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
