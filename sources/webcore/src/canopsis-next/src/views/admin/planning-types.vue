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
    fab-buttons(@create="showCreateTypeModal", @refresh="fetchList")
      span {{ $t('modals.createPbehaviorType.title') }}
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import pbehaviorTypesQueryMixin from '@/mixins/pbehavior/types/query';

import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';
import PlanningTypesList from '@/components/other/pbehavior/types/pbehavior-types-list.vue';

export default {
  components: { FabButtons, PlanningTypesList },
  mixins: [
    rightsTechnicalPbehaviorTypesMixin,
    entitiesPbehaviorTypesMixin,
    pbehaviorTypesQueryMixin,
  ],
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
      this.fetchPbehaviorTypesList({ params: this.getQuery() });
    },

    showCreateTypeModal() {
      this.$modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          action: async (data) => {
            await this.createPbehaviorType({ data });
            await this.fetchList();
          },
        },
      });
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
