<template lang="pug">
  v-card
    v-card-title.green.darken-4.white--text
      v-layout(justify-space-between, align-center)
        h2 {{ config.title }}
        v-btn(@click="hideModal", icon, small)
          v-icon.white--text close
    v-card-text
      v-layout
        v-flex
          chrome(v-model="color")
        v-flex
          compact(v-model="color")
    v-btn(@click="submit") {{ $t('common.submit') }}
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
