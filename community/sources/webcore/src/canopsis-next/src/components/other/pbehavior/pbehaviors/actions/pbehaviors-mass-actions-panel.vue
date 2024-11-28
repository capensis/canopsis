<template>
  <v-layout>
    <c-action-btn
      v-if="removable"
      :tooltip="$t('pbehavior.massRemove')"
      type="delete"
      @click="showRemovePbehaviorsModal"
    />
    <c-action-btn
      v-if="enablable && someOneDisable"
      :tooltip="$t('pbehavior.massEnable')"
      icon="check_circle"
      color="primary"
      @click="showEnablePbehaviorsModal"
    />
    <c-action-btn
      v-if="disablable && someOneEnable"
      :tooltip="$t('pbehavior.massDisable')"
      icon="cancel"
      color="error"
      @click="showDisablePbehaviorsModal"
    />
    <c-db-export-btn :ids="itemsIds" pbehavior />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { mapIds } from '@/helpers/array';
import { pbehaviorToRequest } from '@/helpers/entities/pbehavior/form';

import { useModals } from '@/hooks/modals';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

export default {
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
  setup(props, { emit }) {
    const modals = useModals();
    const {
      fetchPbehaviorsListWithPreviousParams,
      bulkUpdatePbehaviors,
      bulkRemovePbehaviors,
    } = usePbehavior();

    const itemsIds = computed(() => mapIds(props.items));
    const editableItems = computed(() => props.items.filter(({ editable }) => editable));
    const someOneEnable = computed(() => editableItems.value.filter(({ enabled }) => enabled));
    const someOneDisable = computed(() => editableItems.value.filter(({ enabled }) => !enabled));

    /**
     * Clears selected items by emitting clear event
     *
     * @emits clear:items
     */
    const clearItems = () => emit('clear:items');

    /**
     * Clears items and refreshes pbehaviors list
     *
     * @returns {Promise} Result of fetching updated pbehaviors list
     */
    const afterSubmit = async () => {
      clearItems();

      return fetchPbehaviorsListWithPreviousParams();
    };

    /**
     * Shows confirmation modal for enabling selected pbehaviors
     *
     * @returns {Promise} Modal instance
     * @description Opens confirmation modal that:
     * 1. Maps editable items to request format
     * 2. Sets enabled=true for all items
     * 3. Calls bulk update API
     * 4. Clears selection and refreshes list
     */
    const showEnablePbehaviorsModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await bulkUpdatePbehaviors({
            data: editableItems.value.map(item => ({ ...pbehaviorToRequest(item), enabled: true })),
          });

          return afterSubmit();
        },
      },
    });

    /**
     * Shows confirmation modal for disabling selected pbehaviors
     *
     * @returns {Promise} Modal instance
     * @description Opens confirmation modal that:
     * 1. Maps editable items to request format
     * 2. Sets enabled=false for all items
     * 3. Calls bulk update API
     * 4. Clears selection and refreshes list
     */
    const showDisablePbehaviorsModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await bulkUpdatePbehaviors({
            data: editableItems.value.map(item => ({ ...pbehaviorToRequest(item), enabled: false })),
          });

          return afterSubmit();
        },
      },
    });

    /**
     * Shows confirmation modal for removing selected pbehaviors
     *
     * @returns {Promise} Modal instance
     * @description Opens confirmation modal that:
     * 1. Maps selected items to IDs
     * 2. Calls bulk remove API
     * 3. Clears selection and refreshes list
     */
    const showRemovePbehaviorsModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await bulkRemovePbehaviors({ data: mapIds(props.items) });

          return afterSubmit();
        },
      },
    });

    return {
      itemsIds,
      someOneEnable,
      someOneDisable,

      showEnablePbehaviorsModal,
      showDisablePbehaviorsModal,
      showRemovePbehaviorsModal,
    };
  },
};
</script>
