<template>
  <v-layout>
    <c-action-btn
      v-if="updatable"
      :tooltip="pbehavior.editable ? $t('common.edit') : $t('pbehavior.notEditable')"
      :loading="editing"
      type="edit"
      @click="showEditPbehaviorModal"
    />
    <c-action-btn
      v-if="duplicable"
      :loading="duplicating"
      type="duplicate"
      @click="showDuplicatePbehaviorModal"
    />
    <c-action-btn
      v-if="removable"
      type="delete"
      @click="showDeletePbehaviorModal"
    />
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

const { mapActions } = createNamespacedHelpers('pbehavior');

export default {
  inject: ['$system'],
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      editing: false,
      duplicating: false,
    };
  },
  methods: {
    ...mapActions({
      fetchPbehaviorWithoutStore: 'fetchItemWithoutStore',
      removePbehavior: 'remove',
    }),

    refresh() {
      this.$emit('refresh');
    },

    async showEditPbehaviorModal() {
      try {
        this.editing = true;

        const pbehaviorObject = await this.fetchPbehaviorWithoutStore({ id: this.pbehavior._id });

        this.$modals.show({
          name: MODALS.pbehaviorPlanning,
          config: {
            pbehaviors: [pbehaviorObject],
            afterSubmit: this.refresh,
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.editing = false;
      }
    },

    async showDuplicatePbehaviorModal() {
      try {
        this.duplicating = true;

        const pbehaviorObject = await this.fetchPbehaviorWithoutStore({ id: this.pbehavior._id });

        this.$modals.show({
          name: MODALS.pbehaviorPlanning,
          config: {
            pbehaviorsToAdd: [pbehaviorObject],
            afterSubmit: this.refresh,
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.duplicating = false;
      }
    },

    showDeletePbehaviorModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePbehavior({ id: this.pbehavior._id });

            this.refresh();
          },
        },
      });
    },
  },
};
</script>
