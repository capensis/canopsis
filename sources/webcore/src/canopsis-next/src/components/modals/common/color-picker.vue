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
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { Chrome, Compact } from 'vue-color';
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/inner';

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
      const { rgba, hex } = this.color;

      if (this.config.action) {
        const result = this.config.type === 'hex' ? hex : `rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${rgba.a})`;

        await this.config.action(result);
      }

      this.$modals.hide();
    },
  },
};
</script>
