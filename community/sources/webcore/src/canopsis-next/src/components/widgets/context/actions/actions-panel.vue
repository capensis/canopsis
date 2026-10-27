<template>
  <shared-actions-panel
    :actions="actions"
    :inline-count="actions.length"
    ignore-media-query
  />
</template>

<script>
import { compact, pickBy } from 'lodash';

import {
  CONTEXT_ACTIONS_TYPES,
  ENTITY_EXPORT_FILE_NAME_PREFIX,
  ENTITY_TYPES,
  MODALS,
} from '@/constants';

import {
  isPausePropagationAvailableForEntity,
  entityResumeFormToRequest,
  isEntityPaused,
  isEntityPauseInherited,
} from '@/helpers/entities/entity/pause';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';
import { convertObjectToTreeview } from '@/helpers/treeview';

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
    defaultInheritedPbehavior: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    actionsMap() {
      return {
        editEntity: {
          type: CONTEXT_ACTIONS_TYPES.editEntity,
          icon: 'edit',
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
          icon: '$vuetify.icons.pbehavior',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorModal,
        },
        pause: this.pauseAction,
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

    pauseAction() {
      if (!isEntityPaused(this.item)) {
        return {
          type: CONTEXT_ACTIONS_TYPES.pause,
          icon: 'pause',
          title: this.$t('context.actions.titles.pause'),
          method: this.showPauseModal,
        };
      }

      const isInherited = isEntityPauseInherited(this.item);

      return {
        type: CONTEXT_ACTIONS_TYPES.resume,
        icon: 'play_arrow',
        title: isInherited
          ? this.$t('context.actions.titles.inheritedPause', {
            message: String(this.item.pause.message ?? '').trimStart(),
          })
          : this.$t('context.actions.titles.unpause'),
        disabledButton: isInherited,
        method: this.resumeEntityPause,
      };
    },

    actions() {
      const { filteredActionsMap } = this;

      const actions = [
        filteredActionsMap.editEntity,
        filteredActionsMap.pause,
        filteredActionsMap.pbehavior,
        filteredActionsMap.variablesHelp,
      ];

      if (this.item.type === ENTITY_TYPES.service) {
        actions.push(filteredActionsMap.duplicateEntity);
      }

      if (this.item.deletable) {
        actions.push(filteredActionsMap.deleteEntity);
      }

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
              const result = await this.editService({ id: this.item._id, data });

              this.$popups.success({ text: this.$t('modals.createService.success.edit') });

              await this.fetchContextEntitiesListWithPreviousParams();

              return result;
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
            const result = await this.createService({ data });

            this.$popups.success({ text: this.$t('modals.createService.success.duplicate') });

            await this.fetchContextEntitiesListWithPreviousParams();

            return result;
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
          defaultInherited: this.defaultInheritedPbehavior,
          withInherited: true,
        },
      });
    },

    showPauseModal() {
      this.$modals.show({
        name: MODALS.createEntityPauseEvent,
        config: {
          withInherited: isPausePropagationAvailableForEntity(this.item),
          action: async (pause) => {
            try {
              await this.pauseContextEntity({
                data: { ...pause, _id: this.item._id },
              });
              this.$popups.success({ text: this.$t('entity.popups.pauseOperation') });

              return this.fetchContextEntitiesListWithPreviousParams();
            } catch (error) {
              console.error(error);
              this.$popups.error({ text: this.$t('alarm.popups.pauseFailed') });
            }
          },
        },
      });
    },

    async resumeEntityPause() {
      if (isEntityPauseInherited(this.item)) {
        return;
      }

      try {
        await this.resumeContextEntity({
          data: { ...entityResumeFormToRequest(), _id: this.item._id },
        });
        this.$popups.success({ text: this.$t('entity.popups.resumeOperation') });

        return this.fetchContextEntitiesListWithPreviousParams();
      } catch (error) {
        console.error(error);
        this.$popups.error({ text: this.$t('alarm.popups.unpauseFailed') });
      }
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
