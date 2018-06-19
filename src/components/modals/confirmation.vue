<template lang="pug">
  v-card
    v-layout(wrap, justify-space-around)
      v-flex.text-xs-center.py-2.blue.darken-4.white--text(xs12)
        h3 Are you sure ?
      v-flex.text-xs-center.my-2(xs12)
        v-btn.green(small, @click="click(true)") Yes
        v-btn.red(small, @click="click(false)") No
</template>

<script>
import ModalMixin from '@/mixins/modal/modal';

export default {
  name: 'confirmation',
  mixins: [ModalMixin],
  methods: {
    click(choice) {
      this.$emit('confirmation', choice);
      this.hideModal();
    },
  },
};
</script>

