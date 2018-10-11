<template lang="pug">
  v-card
    v-card-title.green.darken-4.white--text
      v-layout(justify-space-between, align-center)
        h2 {{ $t(config.title) }}
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
    return {
      color: {},
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(`
          rgba(${this.color.rgba.r}, ${this.color.rgba.g}, ${this.color.rgba.b}, ${this.color.rgba.a})
        `);
      }
      this.hideModal();
    },
  },
};
</script>
