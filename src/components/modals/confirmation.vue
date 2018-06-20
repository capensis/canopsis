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
import { createNamespacedHelpers } from 'vuex';
import { omit } from 'lodash';
import ModalInnerMixin from '@/mixins/modal/modal-inner';

const { mapActions: contextMapActions } = createNamespacedHelpers('context');

export default {
  name: 'confirmation',
  mixins: [ModalInnerMixin],
  methods: {
    ...contextMapActions({
      removeEntity: 'remove',
    }),
    click(choice) {
      if (choice) {
        let args = { ...this.config.args };
        args = omit(args, ['action']);

        this[this.config.action]({ ...args });
      }

      this.hideModal();
    },
  },
};
</script>

