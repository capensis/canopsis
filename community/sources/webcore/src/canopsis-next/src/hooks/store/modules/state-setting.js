import { storeToRefs } from 'pinia';

import { useStateSettingStore } from '@/stores/modules/state-setting';

/**
 * State setting store facade.
 *
 * @returns {Object}
 */
export const useStateSetting = () => {
  const stateSettingStore = useStateSettingStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(stateSettingStore);

  return {
    stateSettings: items,
    stateSettingsMeta: meta,
    stateSettingsPending: pending,
    stateSettingsFetchingParams: fetchingParams,

    fetchStateSettingsList: stateSettingStore.fetchList,
    fetchStateSettingsListWithPreviousParams: stateSettingStore.fetchListWithPreviousParams,
    createStateSetting: stateSettingStore.create,
    updateStateSetting: stateSettingStore.update,
    removeStateSetting: stateSettingStore.remove,
  };
};
