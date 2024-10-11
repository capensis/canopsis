import Faker from 'faker';
import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreModules } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import {
  CRUD_ACTIONS,
  MODALS,
  PATTERN_TABS,
  PATTERN_TYPES,
  USERS_PERMISSIONS,
} from '@/constants';

import Patterns from '@/views/profile/patterns.vue';

const stubs = {
  'c-page-header': true,
  'c-fab-expand-btn': true,
  'c-action-fab-btn': true,
  'corporate-patterns': true,
  patterns: true,
};

const selectFabExpandButton = wrapper => wrapper.find('c-fab-expand-btn-stub');
const selectPatterns = wrapper => wrapper.find('patterns-stub');
const selectCreateButtons = wrapper => wrapper.findAll('c-fab-expand-btn-stub c-action-fab-btn-stub');
const selectCreateAlarmPatternButton = wrapper => selectCreateButtons(wrapper).at(2);
const selectCreateEntityPatternButton = wrapper => selectCreateButtons(wrapper).at(1);
const selectCreatePbehaviorPatternButton = wrapper => selectCreateButtons(wrapper).at(0);

describe('patterns', () => {
  const $modals = mockModals();

  const fetchCorporatePatternsListWithPreviousParams = jest.fn();
  const fetchPatternsListWithPreviousParams = jest.fn();
  const createPattern = jest.fn();
  const updatePattern = jest.fn();
  const removePattern = jest.fn();
  const bulkRemovePattern = jest.fn();
  const patternModule = {
    name: 'pattern',
    actions: {
      fetchListWithPreviousParams: fetchPatternsListWithPreviousParams,
      create: createPattern,
      update: updatePattern,
      remove: removePattern,
      bulkRemove: bulkRemovePattern,
    },
  };
  const corporatePatternModule = {
    name: 'pattern/corporate',
    actions: {
      fetchListWithPreviousParams: fetchCorporatePatternsListWithPreviousParams,
    },
  };

  const currentUserPermissionsById = jest.fn().mockReturnValue({});
  const authModule = {
    name: 'auth',
    getters: {
      currentUser: {},
      currentUserPermissionsById,
    },
  };
  const store = createMockedStoreModules([
    authModule,
    patternModule,
    corporatePatternModule,
  ]);

  const factory = generateShallowRenderer(Patterns, { stubs });
  const snapshotFactory = generateRenderer(Patterns, { stubs });

  afterEach(() => {
    currentUserPermissionsById.mockClear();
    fetchCorporatePatternsListWithPreviousParams.mockClear();
    fetchPatternsListWithPreviousParams.mockClear();
    createPattern.mockClear();
    updatePattern.mockClear();
    removePattern.mockClear();
    bulkRemovePattern.mockClear();
  });

  it('Patterns refreshed after trigger expand button', async () => {
    const wrapper = factory({
      store,
    });

    const fabButton = selectFabExpandButton(wrapper);

    fabButton.vm.$emit('refresh');

    await flushPromises();

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Corporate patterns refreshed after trigger expand button', async () => {
    const wrapper = factory({
      store,
    });

    const fabButton = selectFabExpandButton(wrapper);

    await wrapper.setData({
      activeTab: PATTERN_TABS.corporatePatterns,
    });

    fabButton.vm.$emit('refresh');

    await flushPromises();

    expect(fetchCorporatePatternsListWithPreviousParams).toBeCalled();
  });

  it('Edit pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const patternsList = selectPatterns(wrapper);

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      type: PATTERN_TYPES.entity,
    };

    patternsList.vm.$emit('edit', pattern);

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.createPattern,
      config: {
        pattern,
        type: PATTERN_TYPES.entity,
        title: 'Edit entity filter',
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const newPattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
    };

    await config.action(newPattern);

    expect(updatePattern).toBeCalledWith(
      expect.any(Object),
      {
        data: newPattern,
        id: pattern._id,
      },
      undefined,
    );

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Edit corporate pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const patternsList = selectPatterns(wrapper);

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      type: PATTERN_TYPES.pbehavior,
      is_corporate: true,
    };

    patternsList.vm.$emit('edit', pattern);

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.createPattern,
      config: {
        pattern,
        type: PATTERN_TYPES.pbehavior,
        title: 'Edit shared pbehavior filter',
        action: expect.any(Function),
      },
    });
  });

  it('Create alarm pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    await flushPromises();

    const createAlarmPatternButton = selectCreateAlarmPatternButton(wrapper);

    createAlarmPatternButton.vm.$emit('click', new Event('click'));

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.createPattern,
      config: {
        pattern: {
          is_corporate: false,
        },
        title: 'Create alarm filter',
        type: PATTERN_TYPES.alarm,
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
    };

    await config.action(pattern);

    expect(createPattern).toBeCalledWith(
      expect.any(Object),
      { data: pattern },
      undefined,
    );

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Create entity pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    await flushPromises();

    const createEntityPatternButton = selectCreateEntityPatternButton(wrapper);

    createEntityPatternButton.vm.$emit('click', new Event('click'));

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.createPattern,
      config: {
        pattern: {
          is_corporate: false,
        },
        title: 'Create entity filter',
        type: PATTERN_TYPES.entity,
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
    };

    await config.action(pattern);

    expect(createPattern).toBeCalledWith(
      expect.any(Object),
      { data: pattern },
      undefined,
    );

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Create pbehavior pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    await flushPromises();

    const createPbehaviorPatternButton = selectCreatePbehaviorPatternButton(wrapper);

    createPbehaviorPatternButton.vm.$emit('click', new Event('click'));

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.createPattern,
      config: {
        pattern: {
          is_corporate: false,
        },
        title: 'Create pbehavior filter',
        type: PATTERN_TYPES.pbehavior,
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
    };

    await config.action(pattern);

    expect(createPattern).toBeCalledWith(
      expect.any(Object),
      { data: pattern },
      undefined,
    );

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Confirmation delete pattern modal showed after trigger delete on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const patternsList = selectPatterns(wrapper);

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      type: PATTERN_TYPES.entity,
    };

    patternsList.vm.$emit('remove', pattern._id);

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(removePattern).toBeCalledWith(
      expect.any(Object),
      { id: pattern._id },
      undefined,
    );
  });

  it('Confirmation delete selected patterns modal showed after trigger delete on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const patternsList = selectPatterns(wrapper);

    const patterns = [
      { _id: Faker.datatype.string() },
      { _id: Faker.datatype.string() },
    ];

    patternsList.vm.$emit('remove-selected', patterns);

    await flushPromises();

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(bulkRemovePattern).toBeCalledWith(
      expect.any(Object),
      { data: patterns },
      undefined,
    );
  });

  it('Renders `patterns` without permissions', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `patterns` with permissions', () => {
    currentUserPermissionsById.mockReturnValueOnce(({
      [USERS_PERMISSIONS.technical.profile.corporatePattern]: {
        actions: [
          CRUD_ACTIONS.create,
          CRUD_ACTIONS.update,
          CRUD_ACTIONS.read,
          CRUD_ACTIONS.delete,
        ],
      },
    }));
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        patternModule,
        corporatePatternModule,
        authModule,
      ]),
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `patterns` corporate tab with permissions', async () => {
    currentUserPermissionsById.mockReturnValueOnce(({
      [USERS_PERMISSIONS.technical.profile.corporatePattern]: {
        actions: [
          CRUD_ACTIONS.create,
          CRUD_ACTIONS.update,
          CRUD_ACTIONS.read,
          CRUD_ACTIONS.delete,
        ],
      },
    }));
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        patternModule,
        corporatePatternModule,
        authModule,
      ]),
    });

    await wrapper.setData({
      activeTab: PATTERN_TABS.corporatePatterns,
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
