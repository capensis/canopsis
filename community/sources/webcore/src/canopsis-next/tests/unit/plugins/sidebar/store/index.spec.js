import { cloneDeep } from 'lodash';
import Faker from 'faker';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import sidebarModule, { types } from '@/plugins/sidebar/store';

const { actions, state: initialState, mutations, getters } = sidebarModule;

describe('Sidebar plugin store module', () => {
  afterEach(() => {
    jest.useRealTimers();
  });

  const createShownState = () => {
    const id = Faker.datatype.uuid();
    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();
    const state = cloneDeep(initialState);

    mutations[types.SHOW](state, { id, name, config });

    return { state, id, name, config };
  };

  it('Mutate state after commit SHOW', () => {
    const state = cloneDeep(initialState);
    const id = Faker.datatype.uuid();
    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();

    mutations[types.SHOW](state, { id, name, config });

    expect(state.allIds).toEqual([id]);
    expect(state.byId[id]).toEqual({
      id,
      name,
      config,
      hidden: false,
      minimized: false,
    });
  });

  it('Mutate state after commit SHOW without config', () => {
    const state = cloneDeep(initialState);
    const id = Faker.datatype.uuid();
    const name = Faker.datatype.string();

    mutations[types.SHOW](state, { id, name });

    expect(state.byId[id].config).toEqual({});
  });

  it('Mutate state after commit HIDE', () => {
    const { state, id } = createShownState();

    mutations[types.HIDE](state, { id });

    expect(state.byId[id].hidden).toBe(true);
  });

  it('Mutate state after commit HIDE_COMPLETED', () => {
    const { state, id } = createShownState();

    mutations[types.HIDE_COMPLETED](state, { id });

    expect(state).toEqual(initialState);
  });

  it('Show sidebar. Action: show', () => {
    const commit = jest.fn();
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();
    const id = Faker.datatype.uuid();
    const payload = { id, name, config };

    actions.show({ commit, state }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.SHOW, payload);
  });

  it('Show sidebar with existing id maximizes. Action: show', () => {
    const commit = jest.fn();
    const { state, id } = createShownState();

    actions.show({ commit, state }, { id, name: 'other', config: {} });

    expect(commit).toHaveBeenCalledWith(types.MAXIMIZE, { id });
  });

  it('Hide sidebar. Action: hide', () => {
    jest.useFakeTimers();
    jest.spyOn(global, 'setTimeout');

    const commit = jest.fn();
    const { state, id } = createShownState();

    actions.hide({ commit, state }, { id });

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.HIDE, { id });

    commit.mockReset();

    expect(setTimeout).toHaveBeenLastCalledWith(
      expect.any(Function),
      VUETIFY_ANIMATION_DELAY,
    );

    jest.runAllTimers();

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.HIDE_COMPLETED, { id });
  });

  it('Hide sidebar skips when id unknown. Action: hide', () => {
    const commit = jest.fn();
    const state = cloneDeep(initialState);

    actions.hide({ commit, state }, { id: 'missing' });

    expect(commit).not.toHaveBeenCalled();
  });

  it('Getter: sidebars', () => {
    const { state, id, name, config } = createShownState();

    expect(getters.sidebars(state)).toEqual([{
      id,
      name,
      config,
      hidden: false,
      minimized: false,
    }]);
  });
});
