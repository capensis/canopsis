import { cloneDeep } from 'lodash';
import Faker from 'faker';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import sidebarModule, { types } from '@/plugins/sidebar/store';

const { actions, state: initialState, mutations, getters } = sidebarModule;

describe('Sidebar plugin store module', () => {
  afterEach(() => {
    jest.useRealTimers();
  });

  it('Mutate state after commit SHOW', () => {
    const show = mutations[types.SHOW];
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();

    show(state, { name, config });

    expect(state).toEqual({ name, config, hidden: false });
  });

  it('Mutate state after commit SHOW without config', () => {
    const show = mutations[types.SHOW];
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();

    show(state, { name });

    expect(state).toEqual({ name, config: {}, hidden: false });
  });

  it('Mutate state after commit HIDE', () => {
    const show = mutations[types.SHOW];
    const hide = mutations[types.HIDE];
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();

    show(state, { name, config });

    hide(state);

    expect(state).toEqual({ name, config, hidden: true });
  });

  it('Mutate state after commit HIDE_COMPLETED', () => {
    const show = mutations[types.SHOW];
    const hideCompleted = mutations[types.HIDE_COMPLETED];
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();

    show(state, { name, config });

    hideCompleted(state);

    expect(state).toEqual(initialState);
  });

  it('Show sidebar. Action: show', () => {
    const commit = jest.fn();
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();
    const payload = { name, config };

    actions.show({ commit, state }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.SHOW, payload);
  });

  it('Hide sidebar. Action: hide', () => {
    jest.useFakeTimers();
    const commit = jest.fn();
    const state = cloneDeep(initialState);

    actions.hide({ commit, state });

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.HIDE);

    commit.mockReset();

    expect(setTimeout).toHaveBeenLastCalledWith(
      expect.any(Function),
      VUETIFY_ANIMATION_DELAY,
    );

    jest.runAllTimers();

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.HIDE_COMPLETED);
  });

  it('Hide sidebar with hidden. Action: hide', () => {
    jest.useFakeTimers();
    const commit = jest.fn();
    const state = { ...initialState, hidden: true };

    actions.hide({ commit, state });

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.HIDE);

    commit.mockReset();

    expect(setTimeout).toHaveBeenLastCalledWith(
      expect.any(Function),
      VUETIFY_ANIMATION_DELAY,
    );

    jest.runAllTimers();

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.HIDE_COMPLETED);
  });

  it('Get sidebar. Getter: sidebar', () => {
    const state = {
      ...initialState,
      config: Faker.helpers.createTransaction(),
    };

    expect(getters.sidebar(state)).toEqual(state);
  });
});
