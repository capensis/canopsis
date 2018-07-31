<template lang='pug'>
  v-card
</template>

<script>
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  mixins: [modalInnerMixin],
  data() {
    return {
    };
  },
  methods: {
  },
};
</script>
