<template>
  <shared-mass-actions-panel :actions="actions" />
</template>

<script>
import { MODALS } from '@/constants';

import { pickIds } from '@/helpers/array';
import { pbehaviorToRequest } from '@/helpers/entities/pbehavior/form';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';

import SharedMassActionsPanel from '@/components/common/actions-panel/mass-actions-panel.vue';

export default {
  components: { SharedMassActionsPanel },
  mixins: [entitiesPbehaviorMixin],
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    removable: {
      type: Boolean,
      default: false,
    },
    enablable: {
      type: Boolean,
      default: false,
    },
    disablable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    editableItems() {
      return this.items.filter(({ editable }) => editable);
    },

    actions() {
      const actions = [];
      const someOneEnable = this.editableItems.some(({ enabled }) => enabled);
      const someOneDisable = this.editableItems.some(({ enabled }) => !enabled);

      if (this.removable) {
        actions.push({
          icon: 'delete',
          iconColor: 'error',
          title: this.$t('pbehavior.massRemove'),
          method: this.showRemovePbehaviorsModal,
        });
      }

      if (this.enablable && someOneDisable) {
        actions.push({
          icon: 'check_circle',
          iconColor: 'primary',
          title: this.$t('pbehavior.massEnable'),
          method: this.showEnablePbehaviorsModal,
        });
      }

      if (this.disablable && someOneEnable) {
        actions.push({
          icon: 'cancel',
          iconColor: 'error',
          title: this.$t('pbehavior.massDisable'),
          method: this.showDisablePbehaviorsModal,
        });
      }

      return actions;
    },
  },
  methods: {
    clearItems() {
      this.$emit('clear:items');
    },

    afterSubmit() {
      this.clearItems();

      return this.fetchPbehaviorsListWithPreviousParams();
    },

    showEnablePbehaviorsModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkUpdatePbehaviors({
              data: this.editableItems.map(item => ({ ...pbehaviorToRequest(item), enabled: true })),
            });

            return this.afterSubmit();
          },
        },
      });
    },

    showDisablePbehaviorsModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkUpdatePbehaviors({
              data: this.editableItems.map(item => ({ ...pbehaviorToRequest(item), enabled: false })),
            });

            return this.afterSubmit();
          },
        },
      });
    },

    showRemovePbehaviorsModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePbehaviors(pickIds(this.items));

            return this.afterSubmit();
          },
        },
      });
    },
  },
};
</script>
