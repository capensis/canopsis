<template lang="pug">
  shared-actions-panel(:actions="actions")
</template>

<script>
import { pickBy } from 'lodash';

import { MODALS, ENTITIES_TYPES, WIDGETS_ACTIONS_TYPES } from '@/constants';

import convertObjectFieldToTreeBranch from '@/helpers/treeview';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import widgetActionsPanelContextMixin from '@/mixins/widget/actions-panel/context';

import SharedActionsPanel from '@/components/other/shared/actions-panel/actions-panel.vue';

/**
 * Component to regroup actions (actions-panel-item) for each entity on the context entities list
 *
 * @module context
 *
 * @prop {Object} item - Item of context entities lists
 * @prop {boolean} [isEditingMode=false] - Is editing mode enable on a view
 */
export default {
  components: { SharedActionsPanel },
  mixins: [
    authMixin,
    modalMixin,
    entitiesWatcherMixin,
    entitiesContextEntityMixin,
    entitiesPbehaviorMixin,
    widgetActionsPanelContextMixin,
  ],
  props: {
    item: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const { context: contextActionsTypes } = WIDGETS_ACTIONS_TYPES;

    return {
      actionsMap: {
        editEntity: {
          type: contextActionsTypes.editEntity,
          icon: 'edit',
          iconClass: 'primary--text',
          title: this.$t('context.actions.titles.editEntity'),
          method: this.showEditEntityModal,
        },
        duplicateEntity: {
          type: contextActionsTypes.duplicateEntity,
          icon: 'file_copy',
          title: this.$t('context.actions.titles.duplicateEntity'),
          method: this.showDuplicateEntityModal,
        },
        deleteEntity: {
          type: contextActionsTypes.deleteEntity,
          icon: 'delete',
          iconClass: 'error--text',
          title: this.$t('context.actions.titles.deleteEntity'),
          method: this.showDeleteEntityModal,
        },
        pbehavior: {
          type: contextActionsTypes.pbehaviorAdd,
          icon: 'pause',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
        variablesHelp: {
          type: contextActionsTypes.variablesHelp,
          icon: 'help',
          title: this.$t('context.actions.titles.variablesHelp'),
          method: this.showVariablesHelpModal,
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

      const actions = [
        filteredActionsMap.editEntity,
        filteredActionsMap.duplicateEntity,
        filteredActionsMap.deleteEntity,
        filteredActionsMap.pbehavior,
      ];

      if (this.isEditingMode) {
        actions.push(filteredActionsMap.variablesHelp);
      }

      return actions.filter(action => !!action);
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
      const parents = [this.item._id];
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

          action: data => this.createPbehavior({
            data,
            parents,
            parentsType,
          }),
        },
      });
    },

    showVariablesHelpModal() {
      const entitiesFields = convertObjectFieldToTreeBranch(this.item, 'entity');
      const variables = [entitiesFields];

      this.showModal({
        name: MODALS.variablesHelp,
        config: {
          variables,
        },
      });
    },
  },
};
</script>
