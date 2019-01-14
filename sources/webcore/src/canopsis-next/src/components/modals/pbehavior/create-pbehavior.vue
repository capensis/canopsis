<template lang="pug">
  pbehavior-form(
  :server-error="serverError",
  :filter="filter",
  @submit="submit",
  @cancel="hideModal",
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';

import PbehaviorForm from '@/components/forms/pbehavior.vue';

const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');

/**
 * Modal to create a pbehavior
 */
export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorForm },
  mixins: [modalInnerItemsMixin],
  data() {
    return {
      serverError: null,
    };
  },
  computed: {
    forEntities() {
      return this.config.itemsIds && this.config.itemsType;
    },

    filter() {
      if (this.forEntities) {
        return {
          _id: { $in: this.items.map(v => v._id) },
        };
      }

      return null;
    },
  },
  methods: {
    ...pbehaviorMapActions({ createPbehavior: 'create' }),

    async submit(data) {
      try {
        this.serverError = null;

        const payload = { data };

        if (this.forEntities) {
          payload.parents = this.items;
          payload.parentsType = this.config.itemsType;
        }

        await this.createPbehavior(payload);

        this.hideModal();
      } catch (err) {
        if (err.description) {
          this.serverError = err.description;
        }
      }
    },
  },
};
</script>
