<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createPause.title') }}
        v-btn(icon, dark, @click="hideModal")
          v-icon close
    v-card-text
      v-textarea(:label="$t('modals.createPause.comment')", v-model="comment")
      v-select(:label="$t('modals.createPause.reason')", v-model="reason", :items="reasons")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(type="submit", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { PAUSE_REASONS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  mixins: [modalInnerMixin],
  data() {
    return {
      reasons: Object.values(PAUSE_REASONS),
      comment: '',
      reason: '',
    };
  },
  methods: {
    async submit() {
      await this.config.action({ comment: this.comment, reason: this.reason });
      this.hideModal();
    },
  },
};
</script>
