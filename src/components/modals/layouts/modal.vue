<template lang="pug">
  v-dialog(v-model="opened", v-bind="dialogProps")
    slot(v-if="opened")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('modal');

export default {
  props: {
    name: {
      type: String,
      required: true,
    },
    dialogProps: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    opened: {
      get() {
        return this.modalComponent === this.name;
      },
      /**
       * setTimeout added for transition
       *
       * @param value
       */
      set(value) {
        if (!value) {
          setTimeout(() => {
            if (this.modalComponent === this.name) {
              this.hideModal();
            }
          }, 300);
        }
      },
    },

    ...mapGetters({
      modalComponent: 'component',
      modalConfig: 'config',
    }),
  },
  methods: {
    ...mapActions({
      hideModal: 'hide',
    }),
  },
};
</script>
