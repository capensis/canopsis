<template lang="pug">
  shared-mass-actions-panel(v-show="items.length", :actions="actions")
</template>

<script>
import { MODALS, WIDGETS_ACTIONS_TYPES } from '@/constants';

import { widgetActionsPanelContextMixin } from '@/mixins/widget/actions-panel/context';

import SharedMassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';
import { pickBy } from 'lodash';

/**
 * Panel regrouping mass actions icons
 *
 * @module context
 *
 * @prop {Array} [itemIds] - Items selected for the mass action
 */
export default {
  components: { SharedMassActionsPanel },
  mixins: [widgetActionsPanelContextMixin],
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    const { context: contextActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actionsMap: {
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
    filteredActionsMap() {
      return pickBy(this.actionsMap, this.actionsAccessFilterHandler);
    },

    actions() {
      const { filteredActionsMap } = this;
      const actions = [filteredActionsMap.pbehavior];
      const everyDeletable = this.items.every(({ deletable }) => deletable);

      if (everyDeletable) {
        actions.unshift(this.filteredActionsMap.deleteEntity);
      }

      return actions.filter(action => !!action);
    },
  },
  methods: {
    showDeleteEntitiesModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const requests = this.items.map(this.removeContextEntityOrService);

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
