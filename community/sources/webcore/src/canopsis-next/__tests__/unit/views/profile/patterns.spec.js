import { mockModals } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MODALS, PATTERN_TABS, PATTERN_TYPES } from '@/constants';

import Patterns from '@/views/profile/patterns-page.vue';

vi.mock('@/hooks/auth', async () => {
  const { computed } = await vi.importActual('vue');

  return {
    useCRUDPermissions: () => ({
      hasCreateAccess: computed(() => true),
      hasReadAccess: computed(() => true),
      hasUpdateAccess: computed(() => true),
      hasDeleteAccess: computed(() => true),
    }),
  };
});

vi.mock('@/hooks/store/modules/pattern', () => ({
  usePattern: () => globalThis.__TEST_PATTERN_HOOK__,
}));

vi.mock('@/hooks/store/modules/pattern-corporate', () => ({
  useCorporatePattern: () => globalThis.__TEST_CORPORATE_PATTERN_HOOK__,
}));

const stubs = {
  'c-page-header': true,
  'c-fab-expand-btn': true,
  'c-action-fab-btn': true,
  'corporate-patterns': true,
  patterns: true,
};

const selectFabExpandButton = wrapper => wrapper.find('c-fab-expand-btn-stub');
const selectCreateButtons = wrapper =>
  wrapper.findAll('c-fab-expand-btn-stub c-action-fab-btn-stub');
const selectCreateAlarmPatternButton = wrapper => selectCreateButtons(wrapper).at(2);
const selectCreateEntityPatternButton = wrapper => selectCreateButtons(wrapper).at(1);
const selectCreatePbehaviorPatternButton = wrapper => selectCreateButtons(wrapper).at(0);

describe('patterns', () => {
  const $modals = mockModals();

  const fetchCorporatePatternsListWithPreviousParams = vi.fn();
  const fetchPatternsListWithPreviousParams = vi.fn();
  const createPattern = vi.fn();
  const updatePattern = vi.fn();
  const removePattern = vi.fn();
  const bulkRemovePattern = vi.fn();
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

  const currentUserPermissionsById = vi.fn().mockReturnValue({});
  const authModule = {
    name: 'auth',
    getters: {
      currentUser: {},
      currentUserPermissionsById,
    },
  };
  const store = createMockedStoreModules([authModule, patternModule, corporatePatternModule]);

  const factory = generateShallowRenderer(Patterns, { stubs });

  beforeEach(() => {
    globalThis.__TEST_PATTERN_HOOK__ = {
      createPattern,
      updatePattern,
      removePattern,
      bulkRemovePatterns: bulkRemovePattern,
      fetchPatternsListWithPreviousParams,
    };
    globalThis.__TEST_CORPORATE_PATTERN_HOOK__ = {
      fetchCorporatePatternsListWithPreviousParams,
    };
  });

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

    fabButton.triggerCustomEvent('refresh');

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

    fabButton.triggerCustomEvent('refresh');

    await flushPromises();

    expect(fetchCorporatePatternsListWithPreviousParams).toBeCalled();
  });

  it('Edit pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      type: PATTERN_TYPES.entity,
    };

    wrapper.vm.showEditPatternModal(pattern);

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
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

    expect(updatePattern).toHaveBeenCalledWith({
      data: newPattern,
      id: pattern._id,
    });

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Edit corporate pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      type: PATTERN_TYPES.pbehavior,
      is_corporate: true,
    };

    wrapper.vm.showEditPatternModal(pattern);

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
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

    createAlarmPatternButton.triggerCustomEvent('click', new Event('click'));

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
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

    expect(createPattern).toHaveBeenCalledWith({ data: pattern });

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Create entity pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    await flushPromises();

    const createEntityPatternButton = selectCreateEntityPatternButton(wrapper);

    createEntityPatternButton.triggerCustomEvent('click', new Event('click'));

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
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

    expect(createPattern).toHaveBeenCalledWith({ data: pattern });

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Create pbehavior pattern modal showed after trigger edit on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    await flushPromises();

    const createPbehaviorPatternButton = selectCreatePbehaviorPatternButton(wrapper);

    createPbehaviorPatternButton.triggerCustomEvent('click', new Event('click'));

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
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

    expect(createPattern).toHaveBeenCalledWith({ data: pattern });

    expect(fetchPatternsListWithPreviousParams).toBeCalled();
  });

  it('Confirmation delete pattern modal showed after trigger delete on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const pattern = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      type: PATTERN_TYPES.entity,
    };

    wrapper.vm.showDeletePatternModal(pattern._id);

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(removePattern).toHaveBeenCalledWith({ id: pattern._id });
  });

  it('Confirmation delete selected patterns modal showed after trigger delete on patterns', async () => {
    const wrapper = factory({ store, mocks: { $modals } });

    const patterns = [{ _id: Faker.datatype.string() }, { _id: Faker.datatype.string() }];

    wrapper.vm.showDeleteSelectedPatternsModal(patterns);

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(bulkRemovePattern).toHaveBeenCalledWith({
      data: patterns.map(({ _id: id }) => ({ _id: id })),
    });
  });
});
