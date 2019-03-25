<template lang="pug">
  shared-mass-actions-panel(v-show="itemsIds.length", :actions="actions")
</template>

<script>
import { MODALS, ENTITIES_TYPES, WIDGETS_ACTIONS_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import popupMixin from '@/mixins/popup';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

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
  mixins: [
    authMixin,
    modalMixin,
    popupMixin,
    entitiesWatcherMixin,
    entitiesPbehaviorMixin,
  ],
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
      const parents = this.itemsIds;
      const parentsType = ENTITIES_TYPES.entity;
      const pbehavior = {
        filter: {
          _id: { $in: [...parents] },
        },
      };

      this.showModal({
        name: MODALS.createPbehavior,
        config: {
          pbehavior,

          action: async (data) => {
            await this.createPbehavior({
              data,
              parents,
              parentsType,
            });

            this.addSuccessPopup({ text: this.$t('success.default') });
          },
        },
      });
    },
  },
};
</script>
