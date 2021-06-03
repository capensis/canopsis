<template lang="pug">
  shared-actions-panel(:actions="actions")
</template>

<script>
import { pickBy } from 'lodash';

import { MODALS, WIDGETS_ACTIONS_TYPES, ENTITY_TYPES } from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';

import { authMixin } from '@/mixins/auth';
import entitiesServiceMixin from '@/mixins/entities/service';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import widgetActionsPanelContextMixin from '@/mixins/widget/actions-panel/context';

import SharedActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

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
    entitiesServiceMixin,
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
          iconColor: 'primary',
          title: this.$t('context.actions.titles.editEntity'),
          method: this.showEditEntityModal,
        },
        duplicateEntity: {
          type: contextActionsTypes.duplicateEntity,
          icon: 'file_copy',
          title: this.$t('context.actions.titles.duplicateEntity'),
          method: this.showDuplicateServiceModal,
        },
        deleteEntity: {
          type: contextActionsTypes.deleteEntity,
          icon: 'delete',
          iconColor: 'error',
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
      ];

      if (this.item.type === ENTITY_TYPES.service) {
        actions.push(filteredActionsMap.duplicateEntity);
      }

      if (this.item.deletable) {
        actions.push(filteredActionsMap.deleteEntity);
      }

      actions.push(filteredActionsMap.pbehavior);

      if (this.isEditingMode) {
        actions.push(filteredActionsMap.variablesHelp);
      }

      return actions.filter(action => !!action);
    },
  },
  methods: {
    async showEditEntityModal() {
      if (this.item.type === ENTITY_TYPES.service) {
        const service = await this.fetchServiceItemWithoutStore({ id: this.item._id });

        this.$modals.show({
          name: MODALS.createService,
          config: {
            item: service,
            title: this.$t('modals.createService.edit.title'),
            action: async (data) => {
              await this.editService({ id: this.item._id, data });

              this.$popups.success({ text: this.$t('modals.createService.success.edit') });

              await this.fetchContextEntitiesListWithPreviousParams();
            },
          },
        });
      } else {
        const basicEntity = await this.fetchBasicContextEntityWithoutStore({ id: this.item._id });

        this.$modals.show({
          name: MODALS.createEntity,
          config: {
            entity: basicEntity,
            title: this.$t('modals.createEntity.edit.title'),
            action: async (data) => {
              await this.updateContextEntityWithPopup({ id: basicEntity._id, data });

              await this.fetchContextEntitiesListWithPreviousParams();
            },
          },
        });
      }
    },

    async showDuplicateServiceModal() {
      const service = await this.fetchServiceItemWithoutStore({ id: this.item._id });

      this.$modals.show({
        name: MODALS.createService,
        config: {
          item: service,
          title: this.$t('modals.createService.duplicate.title'),
          action: async (data) => {
            await this.createService({ data });

            this.$popups.success({ text: this.$t('modals.createService.success.duplicate') });

            await this.fetchContextEntitiesListWithPreviousParams();
          },
        },
      });
    },

    showDeleteEntityModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            if (this.item.type === ENTITY_TYPES.service) {
              await this.removeService({ id: this.item._id });
            } else {
              await this.removeContextEntity({ id: this.item._id });
            }

            await this.fetchContextEntitiesListWithPreviousParams();
          },
        },
      });
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          filter: {
            _id: { $in: [this.item._id] },
          },
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