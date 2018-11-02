<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-card-text
      v-layout
        v-flex
          chrome(v-model="color")
        v-flex
          compact(v-model="color")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { Chrome, Compact } from 'vue-color';
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.colorPicker,
  components: {
    Chrome,
    Compact,
  },
  mixins: [
    modalInnerMixin,
  ],
  data() {
    const { config } = this.modal;
    const data = { color: {} };

    if (config.color) {
      data.color = { hex: config.color };
    }

    return data;
  },
  methods: {
    async submit() {
      const { rgba } = this.color;

      if (this.config.action) {
        await this.config.action(`rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${rgba.a})`);
      }

      this.hideModal();
    },
  },
};
</script>
