<template lang="pug">
  shared-actions-panel(:actions="filteredActions")
</template>

<script>
import { MODALS, ENTITIES_TYPES, WIDGETS_ACTIONS_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import widgetActionsPanelContextMixin from '@/mixins/widget/actions-panel/context';

import SharedActionsPanel from '@/components/other/shared/actions-panel/actions-panel.vue';

/**
 * Component to regroup actions (actions-panel-item) for each entity on the context entities list
 *
 * @module context
 *
 * @prop {Object} item - Item of context entities lists
 */
export default {
  components: { SharedActionsPanel },
  mixins: [
    authMixin,
    modalMixin,
    entitiesWatcherMixin,
    entitiesContextEntityMixin,
    widgetActionsPanelContextMixin,
  ],
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  data() {
    const { context: contextActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actions: [
        {
          type: contextActionsTypes.editEntity,
          icon: 'edit',
          iconClass: 'primary--text',
          title: this.$t('context.actions.titles.editEntity'),
          method: this.showEditEntityModal,
        },
        {
          type: contextActionsTypes.duplicateEntity,
          icon: 'file_copy',
          title: this.$t('context.actions.titles.duplicateEntity'),
          method: this.showDuplicateEntityModal,
        },
        {
          type: contextActionsTypes.deleteEntity,
          icon: 'delete',
          iconClass: 'error--text',
          title: this.$t('context.actions.titles.deleteEntity'),
          method: this.showDeleteEntityModal,
        },
        {
          type: contextActionsTypes.pbehaviorAdd,
          icon: 'pause',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
      ],
    };
  },
  computed: {
    filteredActions() {
      return this.actions.filter(this.actionsAccessFilterHandler);
    },
  },
  methods: {
    showEditEntityModal() {
      if (this.item.type === ENTITIES_TYPES.watcher) {
        this.showModal({
          name: MODALS.createWatcher,
          config: {
            item: this.item,
            title: this.$t('modals.createWatcher.editTitle'),
            action: watcher => this.editWatcherWithPopup(watcher),
          },
        });
      } else {
        this.showModal({
          name: MODALS.createEntity,
          config: {
            item: this.item,
            title: this.$t('modals.createEntity.editTitle'),
            action: entity => this.updateContextEntityWithPopup(entity),
          },
        });
      }
    },

    showDuplicateEntityModal() {
      if (this.item.type === ENTITIES_TYPES.watcher) {
        this.showModal({
          name: MODALS.createWatcher,
          config: {
            item: this.item,
            isDuplicating: true,
            title: this.$t('modals.createWatcher.duplicateTitle'),
            action: watcher => this.duplicateWatcherWithPopup(watcher),
          },
        });
      } else {
        this.showModal({
          name: MODALS.createEntity,
          config: {
            item: this.item,
            isDuplicating: true,
            title: this.$t('modals.createEntity.duplicateTitle'),
            action: entity => this.duplicateContextEntityWithPopup(entity),
          },
        });
      }
    },

    showDeleteEntityModal() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeContextEntity({ id: this.item._id }),
        },
      });
    },

    showAddPbehaviorModal() {
      this.showModal({
        name: MODALS.createPbehavior,
        config: {
          itemsType: ENTITIES_TYPES.entity,
          itemsIds: [this.item._id],
        },
      });
    },
  },
};
</script>
