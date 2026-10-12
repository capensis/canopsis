<template>
  <shared-actions-panel
    v-show="items.length"
    :actions="actions"
    :inline-count="actions.length"
    :small="small"
    ignore-media-query
  />
</template>

<script>
import { CONTEXT_ACTIONS_TYPES, MODALS } from '@/constants';

import { pickIds } from '@/helpers/array';
import {
  isPausePropagationAvailableForEntity,
  entityResumeFormToRequest,
  isEntityPaused,
  isResumeAvailableForEntity,
} from '@/helpers/entities/entity/pause';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import { widgetActionsPanelContextMixin } from '@/mixins/widget/actions-panel/context';

import SharedActionsPanel from '@/components/common/actions-panel/actions-panel.vue';

/**
 * Panel regrouping mass actions icons
 *
 * @module context
 *
 * @prop {Array} [items] - Items selected for the mass action
 */
export default {
  components: { SharedActionsPanel },

  mixins: [widgetActionsPanelContextMixin],

  props: {
    items: {
      type: Array,
      default: () => [],
    },
    defaultInheritedPbehavior: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['clear:items'],

  data() {
    return {
      actionsMap: {
        deleteEntity: {
          type: CONTEXT_ACTIONS_TYPES.deleteEntity,
          icon: 'delete',
          iconColor: 'error',
          title: this.$t('context.actions.titles.deleteEntity'),
          method: this.showDeleteEntitiesModal,
        },
        pbehavior: {
          type: CONTEXT_ACTIONS_TYPES.pbehaviorAdd,
          icon: '$vuetify.icons.pbehavior',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorsModal,
        },
        pause: {
          type: CONTEXT_ACTIONS_TYPES.pause,
          icon: 'pause',
          title: this.$t('context.actions.titles.pause'),
          method: this.showPauseModal,
        },
        resume: {
          type: CONTEXT_ACTIONS_TYPES.resume,
          icon: 'play_arrow',
          title: this.$t('context.actions.titles.unpause'),
          method: this.resumeEntityPauses,
        },
        massEnable: {
          type: CONTEXT_ACTIONS_TYPES.massEnable,
          icon: 'check_circle',
          iconColor: 'primary',
          title: this.$t('context.actions.titles.massEnable'),
          method: this.showEnableEntitiesModal,
        },
        massDisable: {
          type: CONTEXT_ACTIONS_TYPES.massDisable,
          icon: 'cancel',
          iconColor: 'error',
          title: this.$t('context.actions.titles.massDisable'),
          method: this.showDisableEntitiesModal,
        },
      },
    };
  },

  computed: {
    entitiesToPause() {
      return this.items.filter(item => !isEntityPaused(item));
    },

    pausedEntities() {
      return this.items.filter(isEntityPaused);
    },

    entitiesToResume() {
      return this.items.filter(isResumeAvailableForEntity);
    },

    resumeAction() {
      if (!this.pausedEntities.length) {
        return null;
      }

      if (this.entitiesToResume.length) {
        return this.actionsMap.resume;
      }

      return {
        ...this.actionsMap.resume,
        title: this.$t('context.actions.titles.inheritedPauses'),
        disabledButton: true,
      };
    },

    actions() {
      const actions = [];
      const someOneDeletable = this.items.some(({ deletable }) => deletable);
      const someOneEnable = this.items.some(({ enabled }) => enabled);
      const someOneDisable = this.items.some(({ enabled }) => !enabled);

      if (this.entitiesToPause.length) {
        actions.push(this.actionsMap.pause);
      }

      actions.push(this.resumeAction, this.actionsMap.pbehavior);

      if (someOneDisable) {
        actions.push(this.actionsMap.massEnable);
      }

      if (someOneEnable) {
        actions.push(this.actionsMap.massDisable);
      }

      if (someOneDeletable) {
        actions.push(this.actionsMap.deleteEntity);
      }

      return actions.filter(Boolean).filter(this.actionsAccessFilterHandler);
    },
  },

  methods: {
    clearItems() {
      this.$emit('clear:items');
    },

    afterSubmit() {
      this.clearItems();

      return this.fetchContextEntitiesListWithPreviousParams();
    },

    showEnableEntitiesModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkEnableEntities({
              data: pickIds(this.items),
            });

            return this.afterSubmit();
          },
        },
      });
    },

    showDisableEntitiesModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkDisableEntities({
              data: pickIds(this.items),
            });

            return this.afterSubmit();
          },
        },
      });
    },

    showDeleteEntitiesModal() {
      const deletableItems = this.items.filter(({ deletable }) => deletable);

      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          text: this.items.length !== deletableItems.length
            ? this.$t('context.popups.massDeleteWarning')
            : '',
          action: async () => {
            await Promise.all(deletableItems.map(this.removeContextEntityOrService));

            return this.afterSubmit();
          },
        },
      });
    },

    showAddPbehaviorsModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.items.map(({ _id: id }) => id)),
          entities: this.items,
          defaultInherited: this.defaultInheritedPbehavior,
          withInherited: true,
          afterSubmit: this.afterSubmit,
        },
      });
    },

    showPauseModal() {
      this.$modals.show({
        name: MODALS.createEntityPauseEvent,
        config: {
          withInherited: this.entitiesToPause.every(isPausePropagationAvailableForEntity),
          action: async (pause) => {
            try {
              await this.bulkPauseContextEntities({
                data: this.entitiesToPause.map(({ _id: id }) => ({ ...pause, _id: id })),
              });
              this.$popups.success({ text: this.$t('entity.popups.pauseOperation') });

              return this.afterSubmit();
            } catch (error) {
              console.error(error);
              this.$popups.error({ text: this.$t('alarm.popups.pauseFailed') });
            }
          },
        },
      });
    },

    async resumeEntityPauses() {
      if (!this.entitiesToResume.length) {
        return;
      }

      const resume = entityResumeFormToRequest();

      try {
        await this.bulkResumeContextEntities({
          data: this.entitiesToResume.map(({ _id: id }) => ({ ...resume, _id: id })),
        });
        this.$popups.success({ text: this.$t('entity.popups.resumeOperation') });

        return this.afterSubmit();
      } catch (error) {
        console.error(error);
        this.$popups.error({ text: this.$t('alarm.popups.unpauseFailed') });
      }
    },
  },
};
</script>
