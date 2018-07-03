<template lang="pug">
  v-card
    v-layout(wrap, justify-space-around)
      v-flex.text-xs-center.py-2(xs12)
        h3 {{ $t('common.confirmation') }}
      v-flex.text-xs-center.my-2(xs12)
        v-btn.green(small, @click="submit") {{ $t('common.yes') }}
        v-btn.red(small, @click="hideModal") {{ $t('common.no') }}
</template>

<script>
import ModalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: 'confirmation',
  mixins: [ModalInnerMixin],
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action();
      }
      this.hideModal();
    },
  },
};
</script>

