<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('common.confirmation') }}
    v-card-text
      v-layout(wrap, justify-center)
        v-btn.primary(@click="submit") {{ $t('common.yes') }}
        v-btn.error(@click="hideModal") {{ $t('common.no') }}
</template>

<script>
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.confirmation,
  mixins: [modalInnerMixin],
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action();
      }
      this.hideModal();
    },
  },
};
</script>

