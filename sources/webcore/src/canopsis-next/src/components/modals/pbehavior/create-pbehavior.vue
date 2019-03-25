<template lang="pug">
  pbehavior-form(
  :pbehavior="config.pbehavior",
  @submit="submit",
  @cancel="hideModal",
  )
</template>

<script>

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import PbehaviorForm from '@/components/forms/pbehavior.vue';

/**
 * Modal to create a pbehavior
 */
export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorForm },
  mixins: [modalInnerMixin],
  methods: {
    async submit(data) {
      if (this.config.action) {
        await this.config.action(data);
      }

      this.hideModal();
    },
  },
};
</script>
