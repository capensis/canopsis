import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockModals } from '@unit/utils/mock-hooks';

import { MODALS } from '@/constants';

import { uuid } from '@/helpers/uuid';

import FieldFiltersList from '@/components/sidebars/form/fields/filters-list.vue';

jest.mock('@/helpers/uuid');

const stubs = {
  'filters-list': true,
};

const selectFiltersList = wrapper => wrapper.vm.$children[0];

describe('filters-list', () => {
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

  const factory = generateShallowRenderer(FieldFiltersList, {

    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(FieldFiltersList, { stubs });

  it('Filters updated after trigger input event on filters list', () => {
    const wrapper = factory();

    const newFilters = [...filters].reverse();

    selectFiltersList(wrapper).$emit('input', newFilters);

    expect(wrapper).toEmitInput(newFilters);
  });

  it('Filters created after trigger add event on filters list', async () => {
    const widgetId = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        widgetId,
        filters,
      },
    });

    selectFiltersList(wrapper).$emit('add');

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
          entityCountersType: false,
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

    expect(wrapper).toEmitInput([
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

    selectFiltersList(wrapper).$emit('edit', updatedFilter, 1);

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
          entityCountersType: false,
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

    expect(wrapper).toEmitInput([
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

    selectFiltersList(wrapper).$emit('delete', deletedFilter, 1);

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

    expect(wrapper).toEmitInput([
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

    selectFiltersList(wrapper).$emit('delete', deletedFilter, 1);

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

    expect(wrapper).toEmitInput([
      filters[0],
    ]);
  });

  it('Renders `filters` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
