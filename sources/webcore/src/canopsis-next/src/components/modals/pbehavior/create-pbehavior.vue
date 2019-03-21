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

import popupMixin from '@/mixins/popup';
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
  mixins: [popupMixin, modalInnerItemsMixin],
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
      const popups = this.config.popups || {};

      try {
        this.serverError = null;

        const payload = { data };

        if (this.forEntities) {
          payload.parents = this.items;
          payload.parentsType = this.config.itemsType;
        }

        await this.createPbehavior(payload);


        if (popups.success) {
          await this.addSuccessPopup(popups.success);
        }

        this.hideModal();
      } catch (err) {
        if (err.description) {
          this.serverError = err.description;
        }

        if (popups.error) {
          await this.addErrorPopup(popups.error);
        }
      }
    },
  },
};
</script>
