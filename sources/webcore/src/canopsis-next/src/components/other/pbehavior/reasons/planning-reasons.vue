<template lang="pug">
  v-card-text
    planning-reasons-list(
      :pbehavior-reasons="pbehaviorReasons",
      :pending="pbehaviorReasonsPending",
      :totalItems="pbehaviorReasonsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedPbehaviorReasonModal",
      @remove="showRemovePbehaviorReasonModal",
      @edit="showEditPbehaviorReasonModal"
    )
</template>

<script>
import { MODALS } from '@/constants';

import rightsTechnicalPbehaviorReasonsMixin from '@/mixins/rights/technical/pbehavior-reasons';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';
import localQueryMixin from '@/mixins/query-local/query';

import PlanningReasonsList from './pbehavior-reasons-list.vue';

export default {
  components: { PlanningReasonsList },
  mixins: [
    rightsTechnicalPbehaviorReasonsMixin,
    entitiesPbehaviorReasonsMixin,
    localQueryMixin,
  ],
  props: {
    params: {
      type: Object,
      default: () => ({}),
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchPbehaviorReasonsList({ params: this.getQuery() });
    },

    async tryRemovePbehaviorReason(pbehavioReasonId) {
      try {
        await this.removePbehaviorReason({ id: pbehavioReasonId });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    showEditPbehaviorReasonModal(pbehaviorReason) {
      this.$modals.show({
        name: MODALS.createPbehaviorReason,
        config: {
          pbehaviorReason,
          action: async (newPbehaviorReason) => {
            await this.updatePbehaviorReason({
              data: newPbehaviorReason,
              id: pbehaviorReason._id,
            });
            await this.fetchList();
          },
        },
      });
    },

    showRemovePbehaviorReasonModal(pbehaviorReasonId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.tryRemovePbehaviorReason(pbehaviorReasonId);
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedPbehaviorReasonModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.tryRemovePbehaviorReason(id)));

            await this.fetchList();
            this.selected = [];
          },
        },
      });
    },
  },
};
</script>
