<template lang="pug">
  v-btn(@click.stop="deleteEntity", icon, small)
    v-icon delete
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import ModalMixin from '@/mixins/modal/modal';

const { mapActions } = createNamespacedHelpers('context');

export default {
  name: 'context-action-panel',
  mixins: [ModalMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  methods: {
    ...mapActions(['remove']),
    deleteEntity() {
      this.remove({ ids: this.item.props._id });
    },
  },
};
</script>

