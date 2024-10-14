<template>
  <shared-mass-actions-panel
    v-show="items.length"
    :actions="actions"
  />
</template>

<script>
import { MODALS, CONTEXT_ACTIONS_TYPES } from '@/constants';

import { pickIds } from '@/helpers/array';
import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';
import { getPbehaviorNameByEntities } from '@/helpers/entities/pbehavior/form';

import { widgetActionsPanelContextMixin } from '@/mixins/widget/actions-panel/context';

import SharedMassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

/**
 * Panel regrouping mass actions icons
 *
 * @module context
 *
 * @prop {Array} [itemIds] - Items selected for the mass action
 */
export default {
  inject: ['$system'],
  components: { SharedMassActionsPanel },
  mixins: [widgetActionsPanelContextMixin],
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
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
          icon: 'pause',
          title: this.$t('context.actions.titles.pbehavior'),
          method: this.showAddPbehaviorsModal,
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
    actions() {
      const actions = [this.actionsMap.pbehavior];
      const someOneDeletable = this.items.some(({ deletable }) => deletable);
      const someOneEnable = this.items.some(({ enabled }) => enabled);
      const someOneDisable = this.items.some(({ enabled }) => !enabled);

      if (someOneDeletable) {
        actions.unshift(this.actionsMap.deleteEntity);
      }

      if (someOneDisable) {
        actions.push(this.actionsMap.massEnable);
      }

      if (someOneEnable) {
        actions.push(this.actionsMap.massDisable);
      }

      return actions.filter(this.actionsAccessFilterHandler);
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
          defaultName: getPbehaviorNameByEntities(this.items, this.$system.timezone),
          afterSubmit: this.afterSubmit,
        },
      });
    },
  },
};
</script>
