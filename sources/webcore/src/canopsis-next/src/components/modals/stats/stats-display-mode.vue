<template lang="pug">
  v-card(data-test="statsDisplayModeModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('settings.statsNumbers.displayMode') }}
    v-card-text
      stats-display-mode-form(v-model="form")
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(
          data-test="statsDisplayModeCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          data-test="statsDisplayModeSubmitButton",
          @click="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, STATS_DISPLAY_MODE, STATS_DISPLAY_MODE_PARAMETERS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import StatsDisplayModeForm from '@/components/other/stats/stats-display-mode-form.vue';

export default {
  name: MODALS.statsDisplayMode,
  components: { StatsDisplayModeForm },
  mixins: [modalInnerMixin],
  data() {
    const { displayMode } = this.modal.config;
    const defaultDisplayMode = {
      mode: STATS_DISPLAY_MODE.criticity,
      parameters: STATS_DISPLAY_MODE_PARAMETERS,
    };

    return {
      form: cloneDeep(displayMode || defaultDisplayMode),
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.form);
      }

      this.$modals.hide();
    },
  },
};
</script>
