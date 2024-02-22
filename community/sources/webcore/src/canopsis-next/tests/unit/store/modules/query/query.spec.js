import Vue from 'vue';
import { cloneDeep } from 'lodash';
import Faker from 'faker';

import { fakeTimestamp } from '@unit/data/date';
import { mockDateGetTime } from '@unit/utils/mock-hooks';

import SetSeveralPlugin from '@/plugins/set-several';

import queryModule, { types } from '@/store/modules/query';

const { actions, state: initialState, mutations, getters } = queryModule;

Vue.use(SetSeveralPlugin);

describe('Query module', () => {
  const nonceTimestamp = 1326335600000;
  const nowTimestamp = 1386435600000;

  mockDateGetTime(nowTimestamp);

  it('Mutate state after commit UPDATE', () => {
    const update = mutations[types.UPDATE];
    const state = cloneDeep(initialState);

    const id = Faker.datatype.string();
    const query = Faker.helpers.createTransaction();

    update(state, { id, query });

    expect(state).toEqual({
      ...initialState,
      queries: {
        [id]: query,
      },
    });
  });

  it('Mutate state after commit MERGE', () => {
    const merge = mutations[types.MERGE];
    const id = Faker.datatype.string();
    const property = Faker.datatype.string();
    const initialQuery = {
      ...Faker.helpers.createTransaction(),
      property,
    };
    const state = {
      ...initialState,
      queries: {
        [id]: initialQuery,
      },
    };

    const newProperty = Faker.datatype.string();

    merge(state, { id, query: { property: newProperty } });

    expect(state).toEqual({
      ...initialState,
      queries: {
        [id]: {
          ...initialQuery,
          property: newProperty,
        },
      },
    });
  });

  it('Mutate state after commit REMOVE', () => {
    const remove = mutations[types.REMOVE];
    const id = Faker.datatype.string();
    const state = {
      queries: {
        [id]: Faker.helpers.createTransaction(),
      },
      queriesNonces: {
        [id]: fakeTimestamp(),
      },
      lockedQueries: {
        [id]: Faker.helpers.createTransaction(),
      },
    };

    remove(state, { id });

    expect(state).toEqual(initialState);
  });

  it('Mutate state after commit FORCE_UPDATE', () => {
    const forceUpdate = mutations[types.FORCE_UPDATE];
    const id = Faker.datatype.string();
    const state = {
      ...initialState,
      queriesNonces: {
        [id]: nonceTimestamp,
      },
    };

    forceUpdate(state, { id });

    expect(state).toEqual({
      ...initialState,
      queriesNonces: {
        [id]: nowTimestamp,
      },
    });
  });

  it('Mutate state after commit UPDATE_LOCKED', () => {
    const updateLocked = mutations[types.UPDATE_LOCKED];
    const id = Faker.datatype.string();
    const state = {
      ...initialState,
      lockedQueries: {
        [id]: Faker.helpers.createTransaction(),
      },
    };

    const newQuery = Faker.helpers.createTransaction();

    updateLocked(state, { id, query: newQuery });

    expect(state).toEqual({
      ...initialState,
      lockedQueries: {
        [id]: newQuery,
      },
    });
  });

  it('Mutate state after commit REMOVE_LOCKED', () => {
    const removeLocked = mutations[types.REMOVE_LOCKED];
    const id = Faker.datatype.string();
    const state = {
      ...initialState,
      lockedQueries: {
        [id]: Faker.helpers.createTransaction(),
      },
    };

    removeLocked(state, { id });

    expect(state).toEqual(initialState);
  });

  it('Update query. Action: update', () => {
    const commit = jest.fn();

    const payload = {
      id: Faker.datatype.string(),
      query: Faker.helpers.createTransaction(),
    };

    actions.update({ commit }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.UPDATE, payload);
  });

  it('Merge query. Action: merge', () => {
    const commit = jest.fn();

    const payload = {
      id: Faker.datatype.string(),
      query: Faker.helpers.createTransaction(),
    };

    actions.merge({ commit }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.MERGE, payload);
  });

  it('Remove query. Action: remove', () => {
    const commit = jest.fn();

    const payload = {
      id: Faker.datatype.string(),
    };

    actions.remove({ commit }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.REMOVE, payload);
  });

  it('Force update query. Action: forceUpdate', () => {
    const commit = jest.fn();

    const payload = {
      id: Faker.datatype.string(),
    };

    actions.forceUpdate({ commit }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.FORCE_UPDATE, payload);
  });

  it('Update locked query. Action: updateLocked', () => {
    const commit = jest.fn();

    const payload = {
      id: Faker.datatype.string(),
      query: Faker.helpers.createTransaction(),
    };

    actions.updateLocked({ commit }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.UPDATE_LOCKED, payload);
  });

  it('Remove locked query. Action: removeLocked', () => {
    const commit = jest.fn();

    const payload = {
      id: Faker.datatype.string(),
    };

    actions.removeLocked({ commit }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.REMOVE_LOCKED, payload);
  });

  it('Get query by id. Getter: getQueryById', () => {
    const id = Faker.datatype.string();
    const query = Faker.helpers.createTransaction();
    const lockedQuery = Faker.helpers.userCard();
    const state = {
      ...initialState,
      queries: {
        [id]: query,
      },
      lockedQueries: {
        [id]: lockedQuery,
      },
    };

    expect(getters.getQueryById(state)(id)).toEqual({
      ...query,
      ...lockedQuery,
    });
  });

  it('Get query nonce by id. Getter: getQueryNonceById', () => {
    const id = Faker.datatype.string();
    const state = {
      ...initialState,
      queriesNonces: {
        [id]: nonceTimestamp,
      },
    };

    expect(getters.getQueryNonceById(state)(id)).toEqual(nonceTimestamp);
  });

  it('Get query nonce by id without nonce. Getter: getQueryNonceById', () => {
    const id = Faker.datatype.string();

    expect(getters.getQueryNonceById(initialState)(id)).toEqual(0);
  });
});
