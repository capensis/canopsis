<template>
  <shared-actions-panel :actions="actions" />
</template>

<script>
import { compact, pickBy } from 'lodash';

import { MODALS, CONTEXT_ACTIONS_TYPES, ENTITY_TYPES, ENTITY_EXPORT_FILE_NAME_PREFIX } from '@/constants';

import { convertObjectToTreeview } from '@/helpers/treeview';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import { widgetActionsPanelContextMixin } from '@/mixins/widget/actions-panel/context';

import SharedActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

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
    widgetActionsPanelContextMixin,
  ],
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  computed: {
    actionsMap() {
      return {
        editEntity: {
          type: CONTEXT_ACTIONS_TYPES.editEntity,
          icon: 'edit',
          iconColor: 'primary',
          title: this.$t('context.actions.titles.editEntity'),
          method: this.showEditEntityModal,
        },
        duplicateEntity: {
          type: CONTEXT_ACTIONS_TYPES.duplicateEntity,
          icon: 'file_copy',
          title: this.$t('context.actions.titles.duplicateEntity'),
          method: this.showDuplicateServiceModal,
        },
        deleteEntity: {
          type: CONTEXT_ACTIONS_TYPES.deleteEntity,
          icon: 'delete',
          iconColor: 'error',
          title: this.$t('context.actions.titles.deleteEntity'),
          method: this.showDeleteEntityModal,
        },
        pbehavior: {
          type: CONTEXT_ACTIONS_TYPES.pbehaviorAdd,
          icon: 'pause',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
        variablesHelp: {
          type: CONTEXT_ACTIONS_TYPES.variablesHelp,
          icon: 'help',
          title: this.$t('context.actions.titles.variablesHelp'),
          method: this.showVariablesHelpModal,
        },
      };
    },

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

      actions.push(filteredActionsMap.pbehavior, filteredActionsMap.variablesHelp);

      return compact(actions);
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
        const [basicEntity, { impact, depends }] = await Promise.all([
          this.fetchBasicContextEntityWithoutStore({ id: this.item._id }),
          this.fetchContextEntityContextGraphWithoutStore({ id: this.item._id }),
        ]);

        this.$modals.show({
          name: MODALS.createEntity,
          config: {
            entity: {
              ...basicEntity,

              impact,
              depends,
            },
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
            await this.removeContextEntityOrService(this.item);
            await this.fetchContextEntitiesListWithPreviousParams();
          },
        },
      });
    },

    showAddPbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.item._id),
          entities: [this.item],
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
          exportEntity: this.item,
          exportEntityName: `${ENTITY_EXPORT_FILE_NAME_PREFIX}-${this.item._id}`,
        },
      });
    },
  },
};
</script>
