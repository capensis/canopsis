<template lang="pug">
  v-card-text
    planning-dates-exceptions-list(
      :pbehavior-dates-exceptions="pbehaviorDatesExceptions",
      :pending="pbehaviorDatesExceptionsPending",
      :totalItems="pbehaviorDatesExceptionsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedPbehaviorDateExceptionModal",
      @remove="showRemovePbehaviorDateExceptionModal",
      @edit="showEditPbehaviorDateExceptionModal"
    )
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import rightsTechnicalPbehaviorDatesExceptionsMixin from '@/mixins/rights/technical/pbehavior-dates-exceptions';
import entitiesPbehaviorDatesExceptionsMixin from '@/mixins/entities/pbehavior/dates-exceptions';
import pbehaviorQueryMixin from '@/mixins/pbehavior/query';

import PlanningDatesExceptionsList from './pbehavior-dates-exceptions-list.vue';

export default {
  components: { PlanningDatesExceptionsList },
  mixins: [
    rightsTechnicalPbehaviorDatesExceptionsMixin,
    entitiesPbehaviorDatesExceptionsMixin,
    pbehaviorQueryMixin,
  ],
  props: {
    params: {
      type: Object,
      default: () => ({}),
    },
  },
  watch: {
    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
        this.$emit('update:params', this.getQuery());
      }
    },
  },
  mounted() {
    this.fetchList();
    this.$emit('update:params', this.getQuery());
  },
  methods: {
    fetchList() {
      this.fetchPbehaviorDatesExceptionsList({ params: this.getQuery() });
    },

    async tryRemovePbehaviorDateException(pbehavioDateExceptionId) {
      try {
        await this.removePbehaviorDateException({ id: pbehavioDateExceptionId });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    showEditPbehaviorDateExceptionModal(pbehaviorDateException) {
      this.$modals.show({
        name: MODALS.createPbehaviorDateException,
        config: {
          pbehaviorDateException,
          action: async (newPbehaviorDateException) => {
            await this.updatePbehaviorDateException({
              data: newPbehaviorDateException,
              id: pbehaviorDateException._id,
            });
            await this.fetchList();
          },
        },
      });
    },

    showRemovePbehaviorDateExceptionModal(pbehaviorDateExceptionId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.tryRemovePbehaviorDateException(pbehaviorDateExceptionId);
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedPbehaviorDateExceptionModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.tryRemovePbehaviorDateException(id)));

            await this.fetchList();
            this.selected = [];
          },
        },
      });
    },
  },
};
</script>
