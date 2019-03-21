<template lang="pug">
  shared-mass-actions-panel(v-show="itemsIds.length", :actions="actions")
</template>

<script>
import { MODALS, ENTITIES_TYPES, WIDGETS_ACTIONS_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import entitiesWatcherMixin from '@/mixins/entities/watcher';

import SharedMassActionsPanel from '@/components/other/shared/actions-panel/mass-actions-panel.vue';

/**
 * Panel regrouping mass actions icons
 *
 * @module context
 *
 * @prop {Array} [itemIds] - Items selected for the mass action
 */
export default {
  components: { SharedMassActionsPanel },
  mixins: [authMixin, modalMixin, entitiesWatcherMixin],
  props: {
    itemsIds: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    const { context: contextActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actions: [
        {
          type: contextActionsTypes.deleteEntity,
          icon: 'delete',
          iconClass: 'error--text',
          title: this.$t('context.actions.titles.deleteEntity'),
          method: this.showDeleteEntitiesModal,
        },
        {
          type: contextActionsTypes.pbehaviorAdd,
          icon: 'pause',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorsModal,
        },
      ],
    };
  },
  methods: {
    showDeleteEntitiesModal() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => Promise.all(this.itemsIds.map(id => this.removeContextEntity({ id }))),
        },
      });
    },

    showAddPbehaviorsModal() {
      this.showModal({
        name: MODALS.createPbehavior,
        config: {
          itemsType: ENTITIES_TYPES.entity,
          itemsIds: this.itemsIds,
          popups: {
            success: { text: this.$t('success.default') },
            error: { text: this.$t('errors.default') },
          },
        },
      });
    },
  },
};
</script>
