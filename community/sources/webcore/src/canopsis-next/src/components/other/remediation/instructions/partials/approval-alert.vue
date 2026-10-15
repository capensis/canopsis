<template>
  <v-card>
    <v-card-text>
      <v-row class="approval-alert-text" no-gutters>
        <v-col class="approval-alert-text__icon" cols="auto">
          <v-icon color="error">
            notification_important
          </v-icon>
        </v-col>
        <v-col>
          <v-row no-gutters>
            <v-col>
              <strong>{{ userName }}</strong>
              <span class="ml-1">
                {{ statusText }}
              </span>
            </v-col>
          </v-row>
          <v-row class="mt-3" no-gutters>
            <v-col>
              <span class="grey--text text--darken-2">
                {{ $t('common.comment', 1) }}:
              </span>
            </v-col>
          </v-row>
          <v-row class="mt-3" no-gutters>
            <v-col>
              <span>{{ comment }}</span>
            </v-col>
          </v-row>

          <slot />
        </v-col>
      </v-row>
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
