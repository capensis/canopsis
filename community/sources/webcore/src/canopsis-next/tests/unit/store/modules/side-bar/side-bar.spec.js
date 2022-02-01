import { cloneDeep } from 'lodash';
import Faker from 'faker';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import sideBarModule, { types } from '@/store/modules/side-bar';

const { actions, state: initialState, mutations, getters } = sideBarModule;

describe('Side bar module', () => {
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

  it('Show side bar. Action: show', () => {
    const commit = jest.fn();
    const state = cloneDeep(initialState);

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();
    const payload = { name, config };

    actions.show({ commit, state }, payload);

    expect(commit).toHaveBeenCalledTimes(1);
    expect(commit).toHaveBeenCalledWith(types.SHOW, payload);
  });

  it('Show side bar with opened. Action: show', () => {
    const commit = jest.fn();
    const state = { ...initialState, name: 'opened-side-bar' };

    const name = Faker.datatype.string();
    const config = Faker.helpers.createTransaction();
    const payload = { name, config };

    actions.show({ commit, state }, payload);

    expect(commit).not.toHaveBeenCalled();
  });

  it('Hide side bar. Action: hide', () => {
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

    expect(commit).not.toHaveBeenCalled();
  });

  it('Hide side bar with hidden. Action: hide', () => {
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

  it('Get side bar name. Getter: name', () => {
    const state = {
      ...initialState,
      name: Faker.datatype.string(),
    };

    expect(getters.name(state)).toEqual(state.name);
  });

  it('Get side bar config. Getter: config', () => {
    const state = {
      ...initialState,
      config: Faker.helpers.createTransaction(),
    };

    expect(getters.config(state)).toEqual(state.config);
  });

  it('Get side bar hidden. Getter: hidden', () => {
    const state = {
      ...initialState,
      hidden: Faker.datatype.boolean(),
    };

    expect(getters.hidden(state)).toEqual(state.hidden);
  });
});
