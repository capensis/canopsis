<template>
  <v-card>
    <v-card-text>
      <c-layout class="approval-alert-text">
        <div class="approval-alert-text__icon">
          <v-icon color="error">
notification_important
</v-icon>
        </div>
        <c-layout column>
          <c-layout>
            <strong>{{ userName }}</strong>
            <span class="ml-1">
              {{ statusText }}
            </span>
          </c-layout>
          <c-layout class="mt-3">
            <span class="grey--text text--darken-2"> {{ $t('common.comment', 1) }}: </span>
          </c-layout>
          <c-layout class="mt-3">
            <span>{{ comment }}</span>
          </c-layout>

          <slot />
        </c-layout>
      </c-layout>
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
      return this.$t(
        `modals.remediationInstructionApproval.${this.dismissed ? 'dismissed' : 'requested'}`,
      );
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
