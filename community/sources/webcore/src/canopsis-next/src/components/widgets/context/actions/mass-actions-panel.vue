<template lang="pug">
  shared-mass-actions-panel(v-show="items.length", :actions="actions")
</template>

<script>
import { MODALS, CONTEXT_ACTIONS_TYPES, PATTERN_CONDITIONS, ENTITY_PATTERN_FIELDS } from '@/constants';

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
      },
    };
  },
  computed: {
    actions() {
      const actions = [this.actionsMap.pbehavior];
      const someOneDeletable = this.items.some(({ deletable }) => deletable);

      if (someOneDeletable) {
        actions.unshift(this.actionsMap.deleteEntity);
      }

      return actions.filter(this.actionsAccessFilterHandler);
    },
  },
  methods: {
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

            await this.fetchContextEntitiesListWithPreviousParams();
          },
        },
      });
    },

    showAddPbehaviorsModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: [[{
            field: ENTITY_PATTERN_FIELDS.id,
            cond: {
              type: PATTERN_CONDITIONS.isOneOf,
              value: this.items.map(({ _id: id }) => id),
            },
          }]],
        },
      });
    },
  },
};
</script>
