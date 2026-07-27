<template>
  <c-layout gap="1">
    <c-action-btn v-if="removable" type="delete" @click="showDeletePbehaviorModal" />
    <c-db-export-btn :id="pbehavior._id" pbehavior />
  </c-layout>
</template>

<script>
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
  },

  emits: ['refresh'],

  setup(props, { emit }) {
    const modals = useModals();
    const { removePbehavior } = usePbehavior();

    const refresh = () => emit('refresh');

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
      showDeletePbehaviorModal,
    };
  },
};
</script>
