<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      v-btn(@click="showEditPatternModal") {{ $t('modals.eventFilterRule.editPattern') }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createHeartbeat,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        pattern: {},
        expected_interval: '',
      },
    };
  },
  computed: {
    title() {
      if (this.config.heartbeat) {
        return this.$t('modals.createHeartbeat.edit.title');
      }

      return this.$t('modals.createHeartbeat.create.title');
    },
  },
  methods: {
    showEditPatternModal() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          isSimplePattern: true,
          pattern: this.form.pattern,
          action: pattern => this.form.pattern = pattern,
        },
      });
    },

    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.hideModal();
      }
    },
  },
};
</script>
