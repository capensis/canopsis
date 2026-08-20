<template>
  <c-progress-overlay
    :pending="true"
    :opacity="0.8"
    :size="50"
    :width="5"
    class="c-patterns-field__progress-overlay"
    transition
  >
    <template v-if="failedReason" #progress>
      <v-icon size="50" color="error">
        error_outline
      </v-icon>
    </template>
    <v-layout v-if="failedReason" class="gap-4 mt-3" column>
      <span class="error--text text-h6 font-weight-regular">
        {{ failedMessageText || $t('pattern.optimizationFailed', { reason: failedReason }) }}
      </span>
      <v-layout class="gap-2" justify-center>
        <v-btn
          color="primary"
          outlined
          @click="close"
        >
          {{ closeButtonLabel || $t('common.close') }}
        </v-btn>
        <v-btn
          color="primary"
          @click="tryAgain"
        >
          {{ tryAgainButtonLabel || $t('common.tryAgain') }}
        </v-btn>
      </v-layout>
    </v-layout>
    <v-layout
      v-else
      class="gap-4 mt-3"
      justify-center
      column
    >
      <span class="primary--text text-h6 font-weight-regular">
        {{ inProgressText || $t('pattern.optimizationInProgress') }}
      </span>
      <v-layout justify-center>
        <v-btn
          color="primary"
          outlined
          @click="cancel"
        >
          {{ cancelButtonLabel || $t('pattern.cancelOptimization') }}
        </v-btn>
      </v-layout>
    </v-layout>
  </c-progress-overlay>
</template>

<script>
export default {
  props: {
    failedReason: {
      type: String,
      default: '',
    },
    failedMessageText: {
      type: String,
      required: false,
    },
    inProgressText: {
      type: String,
      required: false,
    },
    closeButtonLabel: {
      type: String,
      required: false,
    },
    tryAgainButtonLabel: {
      type: String,
      required: false,
    },
    cancelButtonLabel: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const tryAgain = () => emit('try-again');
    const close = () => emit('close');
    const cancel = () => emit('cancel');

    return {
      tryAgain,
      close,
      cancel,
    };
  },
};
</script>

<style lang="scss">
  .c-patterns-field__progress-overlay .c-progress-overlay__background {
    outline: 3px solid var(--v-background-base);
  }
</style>
