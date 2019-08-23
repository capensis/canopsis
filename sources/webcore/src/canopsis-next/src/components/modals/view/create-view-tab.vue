<template lang="pug">
  v-card(data-test="createTabModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Create tab
    v-card-text
      v-text-field(data-test="tabTitleField", v-model="text")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(data-test="tabSubmitButton", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.createViewTab,
  mixins: [modalInnerMixin],
  data() {
    return {
      text: this.modal.config.text || '',
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.text);
      }

      this.hideModal();
    },
  },
};
</script>

