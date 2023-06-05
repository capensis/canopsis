import Faker from 'faker';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';
import { MODALS } from '@/constants';
import uuid from '@/helpers/uuid';

import Filters from '@/components/sidebars/settings/fields/common/filters.vue';

const localVue = createVueInstance();

jest.mock('@/helpers/uuid');

const stubs = {
  'widget-settings-item': true,
  'filter-selector': true,
  'filters-list': true,
};

const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');
const selectFiltersList = wrapper => wrapper.find('filters-list-stub');

describe('filters', () => {
  const $modals = mockModals();
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

  const factory = generateShallowRenderer(Filters, {
    localVue,
    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(Filters, { localVue, stubs });

  it('Selected filters updated after trigger input on the filter selector field', () => {
    const wrapper = factory();

    selectFilterSelectorField(wrapper).vm.$emit('input', filters[0]);

    expect(wrapper).toEmit('input', filters[0]);
  });

  it('Filters updated after trigger input event on filters list', () => {
    const wrapper = factory();

    const newFilters = [...filters].reverse();

    selectFiltersList(wrapper).vm.$emit('input', newFilters);

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

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createFilter,
        config: {
          title: 'Create filter',
          entityTypes: undefined,
          withAlarm: false,
          withEntity: false,
          withPbehavior: false,
          withServiceWeather: false,
          withTitle: true,
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    const newFilter = {
      name: Faker.datatype.string(),
    };
    const filterId = Faker.datatype.string();

    uuid.mockReturnValueOnce(filterId);

    await config.action(newFilter);

    expect(wrapper).toEmit('update:filters', [
      ...filters,
      {
        ...newFilter,
        _id: filterId,
        widget: widgetId,
        is_private: false,
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

    selectFiltersList(wrapper).vm.$emit('edit', updatedFilter);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createFilter,
        config: {
          filter: updatedFilter,
          title: 'Edit filter',
          entityTypes: undefined,
          withAlarm: false,
          withEntity: false,
          withPbehavior: false,
          withServiceWeather: false,
          withTitle: true,
          action: expect.any(Function),
        },
      },
    );

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

    selectFiltersList(wrapper).vm.$emit('delete', deletedFilter);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(wrapper).toEmit('update:filters', [
      filters[0],
    ]);
  });

  it('Current filter deleted after trigger delete event on filters list', async () => {
    const deletedFilter = filters[1];
    const wrapper = factory({
      propsData: {
        value: deletedFilter._id,
        filters,
      },
    });

    selectFiltersList(wrapper).vm.$emit('delete', deletedFilter);

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];

    await config.action();

    expect(wrapper).toEmit('update:filters', [
      filters[0],
    ]);
    expect(wrapper).toEmit('input', null);
  });

  it('Renders `filters` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `filters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        widgetId: Faker.datatype.string(),
        value: filters[0]._id,
        addable: false,
        editable: false,
        withAlarm: false,
        withEntity: false,
        withPbehavior: false,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
