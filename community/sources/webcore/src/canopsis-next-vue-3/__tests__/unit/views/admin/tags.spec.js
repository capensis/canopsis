import { mockModals } from '@unit/utils/mock-hooks';
import {
  createAuthModule,
  createMockedStoreModules,
  createAlarmTagModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { omit } from 'lodash';

import { MODALS, TAG_TYPES } from '@/constants';

import Tags from '@/views/admin/custom-objects/tags.vue';

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

vi.mock('@/hooks/store/modules/alarm-tag', () => ({
  useAlarmTag: () => globalThis.__TEST_ALARM_TAG_HOOK__,
}));

const stubs = {
  'c-page': true,
  'tags-list': true,
};

describe('tags', () => {
  const $modals = mockModals();

  const { authModule } = createAuthModule();
  const {
    alarmTagModule,

    alarmTags,
    alarmTagsMeta,
    alarmTagsPending,

    fetchAlarmTagsList,
    createAlarmTag,
    updateAlarmTag,
    removeAlarmTag,
    bulkRemoveAlarmTags,
  } = createAlarmTagModule();
  const store = createMockedStoreModules([alarmTagModule, authModule]);

  const factory = generateShallowRenderer(Tags, { stubs });

  beforeEach(() => {
    globalThis.__TEST_ALARM_TAG_HOOK__ = {
      alarmTags,
      alarmTagsMeta,
      alarmTagsPending,
      fetchAlarmTagsList,
      createAlarmTag,
      updateAlarmTag,
      removeAlarmTag,
      bulkRemoveAlarmTag: bulkRemoveAlarmTags,
    };
  });

  test('Tags fetched after mounted', async () => {
    factory({
      store,
    });

    await flushPromises();

    expect(fetchAlarmTagsList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: 10,
        with_flags: true,
      },
    });
  });

  test('Tags re-fetched after trigger refresh button', async () => {
    const wrapper = factory({ store });

    await flushPromises();

    fetchAlarmTagsList.mockClear();

    await wrapper.vm.fetchList();

    await flushPromises();

    expect(fetchAlarmTagsList).toHaveBeenCalledWith({
      params: {
        page: 1,
        limit: 10,
        with_flags: true,
      },
    });
  });

  test('Create tag modal showed after trigger create button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    wrapper.vm.showCreateTagModal();

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createTag,
      config: {
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    const newTag = {
      value: Faker.datatype.string(),
    };

    await modalArguments.config.action(newTag);

    await flushPromises();

    expect(createAlarmTag).toHaveBeenCalledWith({ data: newTag });
    expect(fetchAlarmTagsList).toBeCalled();
  });

  test('Update tag modal showed after trigger edit button', async () => {
    const tag = {
      _id: Faker.datatype.string(),
      type: TAG_TYPES.imported,
    };
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    wrapper.vm.showEditTagModal(tag);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createTag,
      config: {
        tag,
        isImported: true,
        title: 'Edit a tag',
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    const newTag = {
      value: Faker.datatype.string(),
      type: TAG_TYPES.imported,
    };

    await modalArguments.config.action(newTag);

    await flushPromises();

    expect(updateAlarmTag).toHaveBeenCalledWith({
      data: newTag,
      id: tag._id,
    });
    expect(fetchAlarmTagsList).toBeCalled();
  });

  test('Duplicate tag modal showed after trigger duplicate button', async () => {
    const tag = {
      _id: Faker.datatype.string(),
      type: TAG_TYPES.imported,
    };
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    wrapper.vm.showDuplicateTagModal(tag);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createTag,
      config: {
        tag: omit(tag, ['_id']),
        title: 'Duplicate a tag',
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    const newTag = {
      value: Faker.datatype.string(),
      type: TAG_TYPES.imported,
    };

    await modalArguments.config.action(newTag);

    await flushPromises();

    expect(createAlarmTag).toHaveBeenCalledWith({ data: newTag });
    expect(fetchAlarmTagsList).toBeCalled();
  });

  test('Confirmation modal showed after trigger remove tag button', async () => {
    const tag = {
      _id: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    wrapper.vm.showRemoveTagModal(tag._id);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        text: expect.any(String),
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    await modalArguments.config.action();

    await flushPromises();

    expect(removeAlarmTag).toHaveBeenCalledWith({ id: tag._id });
    expect(fetchAlarmTagsList).toBeCalled();
  });

  test('Confirmation modal showed after trigger remove selected tags button', async () => {
    const tag = {
      _id: Faker.datatype.string(),
    };
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    wrapper.vm.showRemoveSelectedTagsModal([tag]);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        text: expect.any(String),
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    await modalArguments.config.action();

    await flushPromises();

    expect(bulkRemoveAlarmTags).toHaveBeenCalledWith({
      data: [{ _id: tag._id }],
    });
    expect(fetchAlarmTagsList).toBeCalled();
  });
});
