import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import CEntityField from '@/components/forms/fields/alarm/c-alarm-tag-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-alarm-action-chip': true,
};

const factory = (options = {}) => shallowMount(CEntityField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CEntityField, {
  localVue,
  stubs,

  ...options,
});
const selectSelectField = wrapper => wrapper.find('v-select-stub');

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
        fetchList: fetchAlarmTags,
      },
    },
  ]);

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

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
