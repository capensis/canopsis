import { defineStore } from 'pinia';
import { ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { useDynamicInfoStore } from '@/stores/modules/dynamic-info';
import { useServiceStore } from '@/stores/modules/service';

export const useInfosStore = defineStore('infos', () => {
  const alarmInfos = ref([]);
  const alarmInfosRules = ref([]);
  const entityInfos = ref([]);
  const pending = ref(false);

  const fetch = async ({ params = { limit: MAX_LIMIT } } = {}) => {
    try {
      pending.value = true;

      const dynamicInfoStore = useDynamicInfoStore();
      const serviceStore = useServiceStore();
      const [
        { data: nextAlarmInfos = [] } = {},
        { data: nextEntityInfos = [] } = {},
        { data: nextAlarmInfosRules = [] } = {},
      ] = await Promise.all([
        dynamicInfoStore.fetchInfosKeysWithoutStore({ params }),
        serviceStore.fetchEntityInfosKeysWithoutStore({ params }),
        dynamicInfoStore.fetchListWithoutStore({ params }),
      ]);

      alarmInfos.value = nextAlarmInfos;
      entityInfos.value = nextEntityInfos;
      alarmInfosRules.value = nextAlarmInfosRules;
    } catch (err) {
      console.error(err);
    } finally {
      pending.value = false;
    }
  };

  return {
    alarmInfos,
    alarmInfosRules,
    entityInfos,
    pending,
    fetch,
  };
});
