import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createSelectInputStub } from '@unit/stubs/input';

import CEntityField from '@/components/forms/fields/alarm/c-alarm-tag-field.vue';
import CLazySearchField from '@/components/forms/fields/c-lazy-search-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field.vue';

const stubs = {
  'c-alarm-action-chip': true,
  'c-lazy-search-field': CLazySearchField,
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-alarm-action-chip': true,
  'c-lazy-search-field': CLazySearchField,
  'c-select-field': CSelectField,
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('c-alarm-tag-field', () => {
  const items = [
    {
      _id: '1',
      value: 'Tag 1',
      color: '#444',
    },
    {
      _id: '2',
      value: 'Tag 2',
      color: '#222',
    },
    {
      _id: '3',
      value: 'Tag 3',
      color: '#000',
    },
  ];
  const fetchAlarmTags = jest.fn().mockReturnValue({
    data: items,
    meta: {
      page_count: items.length,
    },
  });

  const alarmTagsGetter = jest.fn().mockReturnValue(items);
  const pendingGetter = jest.fn().mockReturnValue(false);
  const store = createMockedStoreModules([
    {
      name: 'alarmTag',
      getters: {
        items: alarmTagsGetter,
        pending: pendingGetter,
      },
      actions: {
        fetchListWithoutStore: fetchAlarmTags,
      },
    },
  ]);

  const factory = generateShallowRenderer(CEntityField, { stubs });
  const snapshotFactory = generateRenderer(CEntityField, { stubs: snapshotStubs });

  afterEach(() => {
    fetchAlarmTags.mockClear();
  });

  test('Value changed after trigger the input', async () => {
    const wrapper = factory({
      store,
    });

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input');

    expect(wrapper).toEmit('input');
  });

  test('Renders `c-alarm-tag-field` with default props', async () => {
    const wrapper = snapshotFactory({
      store,
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-alarm-tag-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        value: [items[0].value],
        label: 'Custom label',
        name: 'customName',
        disabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
