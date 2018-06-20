<template lang="pug">
  v-card
    v-layout(wrap, justify-space-around)
      v-flex.text-xs-center.py-2.blue.darken-4.white--text(xs12)
        h3 Are you sure ?
      v-flex.text-xs-center.my-2(xs12)
        v-btn.green(small, @click="submit") Yes
        v-btn.red(small, @click="hideModal") No
</template>

<script>
import ModalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: 'confirmation',
  mixins: [ModalInnerMixin],
  methods: {
    async submit() {
      await this.config.action();

      this.hideModal();
    },
  },
};
</script>

