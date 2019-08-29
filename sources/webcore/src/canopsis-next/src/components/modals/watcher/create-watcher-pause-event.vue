<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createPause.title') }}
        v-btn(icon, dark, @click="hideModal")
          v-icon close
    v-card-text
      v-textarea(
      :label="$t('modals.createPause.comment')",
      v-model="form.comment",
      v-validate="'required'",
      :error-messages="errors.collect('comment')",
      name="comment"
      )
      v-select(
      :label="$t('modals.createPause.reason')",
      v-model="form.reason",
      :items="reasons",
      v-validate="'required'",
      :error-messages="errors.collect('reason')",
      name="reason"
      )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(type="submit", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { PAUSE_REASONS, MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.createWatcherPauseEvent,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      reasons: Object.values(PAUSE_REASONS),
      form: {
        comment: '',
        reason: '',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.hideModal();
      }
    },
  },
};
</script>
