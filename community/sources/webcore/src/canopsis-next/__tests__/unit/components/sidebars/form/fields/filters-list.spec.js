import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { MODALS } from '@/constants';

import { uid } from '@/helpers/uid';

import FieldFiltersList from '@/components/sidebars/form/fields/filters-list.vue';

vi.mock('@/helpers/uid');

const stubs = {
  'filters-list': {
    name: 'FiltersListStub',
    emits: ['update:filters', 'add', 'edit', 'delete'],
    template: '<div class="filters-list-stub" />',
  },
};

const selectFiltersList = wrapper => wrapper.findComponent({ name: 'FiltersListStub' });

describe('filters-list', () => {
  mockModals();
  const $modals = globalThis.__TEST_MODALS__;

  const filters = [
    {
      _id: 'filter-id-1',
      name: 'Filter 1',
    },
    {
      _id: 'filter-id-2',
      name: 'Filter 2',
    },
  ];

  const factory = generateShallowRenderer(FieldFiltersList, {
    stubs,
    mocks: {
      $modals,
    },
  });

  it('Filters updated after trigger input event on filters list', () => {
    const wrapper = factory();

    const newFilters = [...filters].reverse();

    selectFiltersList(wrapper).vm.$emit('update:filters', newFilters);

    expect(wrapper).toEmit('update:filters', newFilters);
  });

  it('Filters created after trigger add event on filters list', async () => {
    const widgetId = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        widgetId,
        filters,
      },
    });

    selectFiltersList(wrapper).vm.$emit('add');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createFilter,
      config: {
        title: 'modals.createFilter.create.title',
        entityTypes: undefined,
        withAlarm: false,
        withEntity: false,
        withPbehavior: false,
        withServiceWeather: false,
        entityCountersType: false,
        widgetType: undefined,
        withTitle: true,
        withoutLink: true,
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    const newFilter = {
      name: Faker.datatype.string(),
    };
    const filterId = Faker.datatype.string();

    uid.mockReturnValueOnce(filterId);

    await config.action(newFilter);

    expect(wrapper).toEmit('update:filters', [
      ...filters,
      {
        ...newFilter,
        _id: filterId,
        widget: widgetId,
        is_user_preference: false,
      },
    ]);
  });

  it('Filters edited after trigger edit event on filters list', async () => {
    const widgetId = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        widgetId,
        filters,
      },
    });

    const updatedFilter = filters[1];

    selectFiltersList(wrapper).vm.$emit('edit', updatedFilter, 1);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createFilter,
      config: {
        filter: updatedFilter,
        title: 'modals.createFilter.edit.title',
        entityTypes: undefined,
        withAlarm: false,
        withEntity: false,
        withPbehavior: false,
        withServiceWeather: false,
        entityCountersType: false,
        widgetType: undefined,
        withTitle: true,
        withoutLink: true,
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const newData = {
      name: Faker.datatype.string(),
    };

    await config.action(newData);

    expect(wrapper).toEmit('update:filters', [
      filters[0],
      {
        ...updatedFilter,
        ...newData,
        widget: widgetId,
      },
    ]);
  });

  it('Filter deleted after trigger delete event on filters list', async () => {
    const wrapper = factory({
      propsData: {
        filters,
      },
    });

    const deletedFilter = filters[1];

    selectFiltersList(wrapper).vm.$emit('delete', deletedFilter, 1);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(wrapper).toEmit('update:filters', [filters[0]]);
  });

  it('Current filter deleted after trigger delete event on filters list', async () => {
    const deletedFilter = filters[1];
    const wrapper = factory({
      propsData: {
        filters,
      },
    });

    selectFiltersList(wrapper).vm.$emit('delete', deletedFilter, 1);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(wrapper).toEmit('update:filters', [filters[0]]);
  });
});
