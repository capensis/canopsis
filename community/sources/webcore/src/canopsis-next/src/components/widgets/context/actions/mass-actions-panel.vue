<template lang="pug">
  shared-mass-actions-panel(v-show="items.length", :actions="filteredActions")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITY_TYPES, MODALS, WIDGETS_ACTIONS_TYPES } from '@/constants';

import SharedMassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

const { mapActions: mapContextEntityActions } = createNamespacedHelpers('entity');
const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

/**
 * Panel regrouping mass actions icons
 *
 * @module context
 *
 * @prop {Array} [itemIds] - Items selected for the mass action
 */
export default {
  components: { SharedMassActionsPanel },

  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    const { context: contextActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actions: {
        deleteEntity: {
          type: contextActionsTypes.deleteEntity,
          icon: 'delete',
          iconColor: 'error',
          title: this.$t('context.actions.titles.deleteEntity'),
          method: this.showDeleteEntitiesModal,
        },
        pbehavior: {
          type: contextActionsTypes.pbehaviorAdd,
          icon: 'pause',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorsModal,
        },
      },
    };
  },
  computed: {
    filteredActions() {
      const actions = [this.actions.pbehavior];
      const everyDeletable = this.items.every(({ deletable }) => deletable);

      if (everyDeletable) {
        actions.unshift(this.actions.deleteEntity);
      }

      return actions;
    },
  },
  methods: {
    ...mapContextEntityActions({
      removeContextEntity: 'remove',
      fetchContextEntitiesListWithPreviousParams: 'fetchListWithPreviousParams',
    }),

    ...mapServiceActions({
      removeService: 'remove',
    }),

    showDeleteEntitiesModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const requests = this.items.map((item) => {
              if (item.type === ENTITY_TYPES.service) {
                return this.removeService({ id: item._id });
              }

              return this.removeContextEntity({ id: item._id });
            });

            await Promise.all(requests);

            await this.fetchContextEntitiesListWithPreviousParams();
          },
        },
      });
    },

    showAddPbehaviorsModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          filter: {
            _id: { $in: this.items.map(({ _id: id }) => id) },
          },
        },
      });
    },
  },
};
</script>
