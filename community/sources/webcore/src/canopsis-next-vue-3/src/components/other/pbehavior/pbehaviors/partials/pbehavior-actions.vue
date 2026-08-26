<template>
  <c-layout gap="1">
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
    <c-action-btn v-if="removable" type="delete" @click="showDeletePbehaviorModal" />
    <c-db-export-btn :id="pbehavior._id" pbehavior />
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

export default {
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
    withInherited: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['refresh'],

  setup(props, { emit }) {
    const modals = useModals();
    const { fetchPbehaviorWithoutStore, removePbehavior } = usePbehavior();
    const editing = ref(false);
    const duplicating = ref(false);

    const refresh = () => emit('refresh');

    const showEditPbehaviorModal = async () => {
      try {
        editing.value = true;

        const pbehaviorObject = await fetchPbehaviorWithoutStore({ id: props.pbehavior._id });

        modals.show({
          name: MODALS.pbehaviorPlanning,
          config: {
            pbehaviors: [pbehaviorObject],
            withInherited: props.withInherited,
            afterSubmit: refresh,
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        editing.value = false;
      }
    };

    const showDuplicatePbehaviorModal = async () => {
      try {
        duplicating.value = true;

        const pbehaviorObject = await fetchPbehaviorWithoutStore({ id: props.pbehavior._id });

        modals.show({
          name: MODALS.pbehaviorPlanning,
          config: {
            pbehaviorsToAdd: [pbehaviorObject],
            withInherited: props.withInherited,
            afterSubmit: refresh,
          },
        });
      } catch (err) {
        console.error(err);
      } finally {
        duplicating.value = false;
      }
    };

    const showDeletePbehaviorModal = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removePbehavior({ id: props.pbehavior._id });

            refresh();
          },
        },
      });
    };

    return {
      editing,
      duplicating,
      showEditPbehaviorModal,
      showDuplicatePbehaviorModal,
      showDeletePbehaviorModal,
    };
  },
};
</script>
