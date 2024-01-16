<template>
  <v-card-text>
    <planning-exceptions-list
      :pbehavior-exceptions="pbehaviorExceptions"
      :pending="pbehaviorExceptionsPending"
      :total-items="pbehaviorExceptionsMeta.total_count"
      :options.sync="options"
      @remove-selected="showRemoveSelectedPbehaviorExceptionModal"
      @remove="showRemovePbehaviorExceptionModal"
      @edit="showEditPbehaviorExceptionModal"
    />
  </v-card-text>
</template>

<script>
import { MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesPbehaviorExceptionMixin } from '@/mixins/entities/pbehavior/exceptions';
import { permissionsTechnicalPbehaviorExceptionsMixin } from '@/mixins/permissions/technical/pbehavior-exceptions';

import PlanningExceptionsList from './pbehavior-exceptions-list.vue';

export default {
  components: { PlanningExceptionsList },
  mixins: [
    localQueryMixin,
    entitiesPbehaviorExceptionMixin,
    permissionsTechnicalPbehaviorExceptionsMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;

      this.fetchPbehaviorExceptionsList({ params });
    },

    async tryRemovePbehaviorException(pbehavioExceptionId) {
      try {
        await this.removePbehaviorException({ id: pbehavioExceptionId });
      } catch (err) {
        console.error(err);

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
          },
        },
      });
    },
  },
};
</script>
