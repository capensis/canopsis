<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createHeartbeat.create.title') }}
    v-card-text
      heartbeat-form(v-model="form")
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import HeartbeatForm from '@/components/other/heartbeat/form/heartbeat-form.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createHeartbeat,
  $_veeValidate: {
    validator: 'new',
  },
  components: { HeartbeatForm },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        pattern: {},
        periodValue: '',
        periodUnit: '',
      },
    };
  },
  methods: {
    async submit() {
      try {
        const isValid = await this.$validator.validateAll();

        if (isValid) {
          const { pattern, periodValue, periodUnit } = this.form;
          const data = {
            pattern,
            expected_interval: `${periodValue}${periodUnit}`,
          };

          if (this.config.action) {
            await this.config.action(data);
          }

          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: err.description });
      }
    },
  },
};
</script>
