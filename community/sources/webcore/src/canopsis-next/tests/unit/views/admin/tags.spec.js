import Faker from 'faker';
import { omit } from 'lodash';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createAuthModule, createMockedStoreModules, createAlarmTagModule } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';

import { CRUD_ACTIONS, MODALS, TAG_TYPES, USERS_PERMISSIONS } from '@/constants';

import Tags from '@/views/admin/tags.vue';

const stubs = {
  'c-page': true,
  'tags-list': true,
};

const selectPageNode = wrapper => wrapper.vm.$children[0];
const selectTagsList = wrapper => wrapper.find('tags-list-stub');

describe('tags', () => {
  const $modals = mockModals();

  const { authModule, currentUserPermissionsById } = createAuthModule();
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
  const store = createMockedStoreModules([
    alarmTagModule,
    authModule,
  ]);

  const factory = generateShallowRenderer(Tags, { stubs });
  const snapshotFactory = generateRenderer(Tags, { stubs });

  test('Tags fetched after mounted', async () => {
    factory({
      store,
    });

    await flushPromises();

    expect(fetchAlarmTagsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          page: 1,
          limit: 10,
          with_flags: true,
        },
      },
      undefined,
    );
  });

  test('Tags re-fetched after trigger refresh button', async () => {
    const wrapper = factory({ store });

    await flushPromises();

    fetchAlarmTagsList.mockClear();

    selectPageNode(wrapper).$emit('refresh');

    await flushPromises();

    expect(fetchAlarmTagsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          page: 1,
          limit: 10,
          with_flags: true,
        },
      },
      undefined,
    );
  });

  test('Create tag modal showed after trigger create button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    selectPageNode(wrapper).$emit('create');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createTag,
        config: {
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    const newTag = {
      value: Faker.datatype.string(),
    };

    modalArguments.config.action(newTag);

    await flushPromises();

    expect(createAlarmTag).toBeCalledWith(
      expect.any(Object),
      {
        data: newTag,
      },
      undefined,
    );
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

    selectTagsList(wrapper).triggerCustomEvent('edit', tag);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createTag,
        config: {
          tag,
          isImported: true,
          title: 'Edit a tag',
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    const newTag = {
      value: Faker.datatype.string(),
      type: TAG_TYPES.imported,
    };

    modalArguments.config.action(newTag);

    await flushPromises();

    expect(updateAlarmTag).toBeCalledWith(
      expect.any(Object),
      {
        data: newTag,
        id: tag._id,
      },
      undefined,
    );
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

    selectTagsList(wrapper).triggerCustomEvent('duplicate', tag);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createTag,
        config: {
          tag: omit(tag, ['_id']),
          title: 'Duplicate a tag',
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    const newTag = {
      value: Faker.datatype.string(),
      type: TAG_TYPES.imported,
    };

    modalArguments.config.action(newTag);

    await flushPromises();

    expect(createAlarmTag).toBeCalledWith(
      expect.any(Object),
      {
        data: newTag,
      },
      undefined,
    );
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

    await selectTagsList(wrapper).triggerCustomEvent('remove', tag._id);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          text: expect.any(String),
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    await flushPromises();

    expect(removeAlarmTag).toBeCalledWith(
      expect.any(Object),
      {
        id: tag._id,
      },
      undefined,
    );
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

    await selectTagsList(wrapper).triggerCustomEvent('remove-selected', [tag]);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          text: expect.any(String),
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    await flushPromises();

    expect(bulkRemoveAlarmTags).toBeCalledWith(
      expect.any(Object),
      {
        data: [{ _id: tag._id }],
      },
      undefined,
    );
    expect(fetchAlarmTagsList).toBeCalled();
  });

  test('Renders `tags` without permissions', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `tags` with permissions', () => {
    currentUserPermissionsById.mockReturnValueOnce(({
      [USERS_PERMISSIONS.technical.tag]: {
        actions: [
          CRUD_ACTIONS.create,
          CRUD_ACTIONS.update,
          CRUD_ACTIONS.read,
          CRUD_ACTIONS.delete,
        ],
      },
    }));

    alarmTags.mockReturnValueOnce([
      { _id: 'first-tag' },
      { _id: 'second-tag' },
    ]);
    alarmTagsMeta.mockReturnValueOnce({
      total_count: 2,
    });
    alarmTagsPending.mockReturnValueOnce(true);
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        alarmTagModule,
        authModule,
      ]),
    });

    expect(wrapper).toMatchSnapshot();
  });
});
