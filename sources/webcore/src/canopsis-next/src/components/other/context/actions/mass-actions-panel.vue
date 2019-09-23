<template lang="pug">
  shared-mass-actions-panel(v-show="itemsIds.length", :actions="actions")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, WIDGETS_ACTIONS_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import SharedMassActionsPanel from '@/components/other/shared/actions-panel/mass-actions-panel.vue';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

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
    entitiesWatcherMixin,
    entitiesContextEntityMixin,
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
  computed: {
    ...entitiesMapGetters({
      getEntitiesList: 'getList',
    }),
  },
  methods: {
    showDeleteEntitiesModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => Promise.all(this.itemsIds.map(id => this.removeContextEntity({ id }))),
        },
      });
    },

    showAddPbehaviorsModal() {
      this.$modals.show({
        name: MODALS.createPbehavior,
        config: {
          pbehavior: {
            filter: {
              _id: { $in: this.itemsIds },
            },
          },
          action: async (data) => {
            await this.createPbehavior({ data });

            this.$popups.addSuccess({ text: this.$t('success.default') });
          },
        },
      });
    },
  },
};
</script>
