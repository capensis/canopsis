import { merge, omit } from 'lodash';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useQueryStore = defineStore('query', () => {
  const queries = ref({});
  const queriesNonces = ref({});
  const lockedQueries = ref({});

  const getQueryById = computed(
    () => id => ({
      ...queries.value[id],
      ...lockedQueries.value[id],
    }),
  );
  const getQueryNonceById = computed(() => id => queriesNonces.value[id] || 0);

  const update = ({ id, query } = {}) => {
    queries.value = {
      ...queries.value,
      [id]: query,
    };
  };

  const mergeQuery = ({ id, query } = {}) => {
    queries.value = {
      ...queries.value,
      [id]: merge({}, queries.value[id], query),
    };
  };

  const remove = ({ id } = {}) => {
    queries.value = omit(queries.value, [id]);
    queriesNonces.value = omit(queriesNonces.value, [id]);
    lockedQueries.value = omit(lockedQueries.value, [id]);
  };

  const forceUpdate = ({ id } = {}) => {
    queriesNonces.value = {
      ...queriesNonces.value,
      [id]: Date.now(),
    };
  };

  const updateLocked = ({ id, query } = {}) => {
    lockedQueries.value = {
      ...lockedQueries.value,
      [id]: query,
    };
  };

  const removeLocked = ({ id } = {}) => {
    lockedQueries.value = omit(lockedQueries.value, [id]);
  };

  return {
    queries,
    queriesNonces,
    lockedQueries,

    getQueryById,
    getQueryNonceById,

    update,
    merge: mergeQuery,
    remove,
    forceUpdate,
    updateLocked,
    removeLocked,
  };
});
