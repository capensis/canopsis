<template lang="pug">
  v-card(data-test="confirmationModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('common.confirmation') }}
    v-card-text
      v-layout(wrap, justify-center)
        v-btn.primary(
        @click.prevent="submit",
        data-test="confirmButton",
        :loading="submitting",
        :disabled="submitting"
        ) {{ $t('common.yes') }}
        v-btn.error(@click="hideModal") {{ $t('common.no') }}
</template>

<script>
import modalInnerMixin from '@/mixins/modal/inner';
import { MODALS } from '@/constants';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.confirmation,
  mixins: [modalInnerMixin],
  data() {
    return {
      submitting: false,
    };
  },
  methods: {
    async submit() {
      this.submitting = true;
      if (this.config.action) {
        await this.config.action();
      }
      this.hideModal();
      this.submitting = false;
    },
  },
};
</script>

