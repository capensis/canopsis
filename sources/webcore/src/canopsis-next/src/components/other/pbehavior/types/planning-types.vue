<template lang="pug">
  v-card-text
    planning-types-list(
      :pbehavior-types="pbehaviorTypes",
      :pending="pbehaviorTypesPending",
      :totalItems="pbehaviorTypesMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedPbehaviorTypeModal",
      @remove="showRemovePbehaviorTypeModal",
      @edit="showEditPbehaviorTypeModal"
    )
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS, PLANNING_TABS } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import pbehaviorQueryMixin from '@/mixins/pbehavior/query';

import PlanningTypesList from '@/components/other/pbehavior/types/pbehavior-types-list.vue';

export default {
  components: { PlanningTypesList },
  mixins: [
    rightsTechnicalPbehaviorTypesMixin,
    entitiesPbehaviorTypesMixin,
    pbehaviorQueryMixin,
  ],
  props: {
    params: {
      type: Object,
      default: () => ({}),
    },
    queryId: {
      type: String,
      default: PLANNING_TABS.types,
    },
  },
  watch: {
    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchPbehaviorTypesList({ params: this.query });
    },

    async tryRemovePbehaviorType(pbehaviorTypeId) {
      try {
        await this.removePbehaviorType({ id: pbehaviorTypeId });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    showEditPbehaviorTypeModal(pbehaviorType) {
      this.$modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          pbehaviorType,
          action: async (newPbehaviorType) => {
            await this.updatePbehaviorType({
              data: newPbehaviorType,
              id: pbehaviorType._id,
            });
            await this.fetchList();
          },
        },
      });
    },

    showRemovePbehaviorTypeModal(pbehaviorTypeId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.tryRemovePbehaviorType(pbehaviorTypeId);
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedPbehaviorTypeModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.tryRemovePbehaviorType(id)));

            await this.fetchList();
            this.selected = [];
          },
        },
      });
    },
  },
};
</script>
