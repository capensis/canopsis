<template lang="pug">
  v-card-text
    planning-exceptions-list(
      :pbehaviorExceptions="pbehaviorExceptions",
      :pending="pbehaviorExceptionsPending",
      :totalItems="pbehaviorExceptionsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedPbehaviorExceptionModal",
      @remove="showRemovePbehaviorExceptionModal",
      @edit="showEditPbehaviorExceptionModal"
    )
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS, PLANNING_TABS } from '@/constants';

import rightsTechnicalPbehaviorExceptionsMixin from '@/mixins/rights/technical/pbehavior-exceptions';
import entitiesPbehaviorExceptionsMixin from '@/mixins/entities/pbehavior/exceptions';
import pbehaviorQueryMixin from '@/mixins/pbehavior/query';

import PlanningExceptionsList from './pbehavior-exceptions-list.vue';

export default {
  components: { PlanningExceptionsList },
  mixins: [
    rightsTechnicalPbehaviorExceptionsMixin,
    entitiesPbehaviorExceptionsMixin,
    pbehaviorQueryMixin,
  ],
  props: {
    queryId: {
      type: String,
      default: PLANNING_TABS.exceptions,
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
      this.fetchPbehaviorExceptionsList({ params: this.query });
    },

    async tryRemovePbehaviorException(pbehavioExceptionId) {
      try {
        await this.removePbehaviorException({ id: pbehavioExceptionId });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    showEditPbehaviorExceptionModal(pbehaviorException) {
      this.$modals.show({
        name: MODALS.createPbehaviorException,
        config: {
          pbehaviorException,
          action: async (newPbehaviorException) => {
            await this.updatePbehaviorException({
              data: newPbehaviorException,
              id: pbehaviorException._id,
            });
            await this.fetchList();
          },
        },
      });
    },

    showRemovePbehaviorExceptionModal(pbehaviorExceptionId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.tryRemovePbehaviorException(pbehaviorExceptionId);
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedPbehaviorExceptionModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.tryRemovePbehaviorException(id)));

            await this.fetchList();
            this.selected = [];
          },
        },
      });
    },
  },
};
</script>
