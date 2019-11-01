<template lang="pug">
  shared-actions-panel(:actions="actions")
</template>

<script>
import { pickBy } from 'lodash';

import { MODALS, ENTITIES_TYPES, WIDGETS_ACTIONS_TYPES } from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';

import authMixin from '@/mixins/auth';
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
        this.$modals.show({
          name: MODALS.createWatcher,
          config: {
            item: this.item,
            title: this.$t('modals.createWatcher.editTitle'),
            action: watcher => this.editWatcherWithPopup(watcher),
          },
        });
      } else {
        this.$modals.show({
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
        this.$modals.show({
          name: MODALS.createWatcher,
          config: {
            item: this.item,
            isDuplicating: true,
            title: this.$t('modals.createWatcher.duplicateTitle'),
            action: watcher => this.duplicateWatcherWithPopup(watcher),
          },
        });
      } else {
        this.$modals.show({
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
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeContextEntity({ id: this.item._id }),
        },
      });
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.createPbehavior,
        config: {
          pbehavior: {
            filter: {
              _id: { $in: [this.item._id] },
            },
          },
          action: data => this.createPbehavior({ data }),
        },
      });
    },

    showVariablesHelpModal() {
      const entitiesFields = convertObjectToTreeview(this.item, 'entity');
      const variables = [entitiesFields];

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          variables,
        },
      });
    },
  },
};
</script>
