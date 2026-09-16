import AxiosMockAdapter from 'axios-mock-adapter';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals, mockPopups } from '@unit/utils/mock-hooks';
import {
  createMockedStoreModules,
  createUserPreferenceModule,
  createWidgetModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { API_ROUTES } from '@/config';
import { MODALS } from '@/constants';

import ClickOutside from '@/services/click-outside';
import request from '@/services/request';

import FiltersList from '@/components/modals/common/filters-list.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-progress-overlay': true,
  'filters-list-component': {
    props: ['filters', 'pending', 'addable', 'editable'],
    emits: ['update:filters', 'add', 'edit', 'delete'],
    template: '<div class="filters-list-component-stub" />',
  },
};

const selectFiltersList = wrapper => wrapper.find('.filters-list-component-stub');

describe('filters-list', () => {
  const $modals = mockModals();
  const $popups = mockPopups();

  const { userPreferenceModule, fetchUserPreference, getUserPreferenceByWidgetId } =
    createUserPreferenceModule();
  const {
    widgetModule,
    createWidgetFilter,
    updateWidgetFilter,
    removeWidgetFilter,
    updateWidgetFiltersPositions,
  } = createWidgetModule();
  const store = createMockedStoreModules([userPreferenceModule, widgetModule]);

  const axiosMockAdapter = new AxiosMockAdapter(request);

  axiosMockAdapter.onAny().passThrough();

  beforeEach(() => {
    axiosMockAdapter.reset();
    axiosMockAdapter.onGet(API_ROUTES.patternFields.widgetFilter).reply(200, {
      entity_pattern: [
        { name: 'name', enabled: true, alias: false },
        { name: 'category', enabled: true, alias: false },
        { name: 'component', enabled: true, alias: false },
      ],
      alarm_pattern: [],
      event_pattern: [],
      pbehavior_pattern: [],
      weather_service_pattern: [],
    });
  });

  const widgetId = 'widget-id';
  const modal = {
    id: 'test-modal-id',
    config: {
      widget: { _id: widgetId },
      addable: true,
      editable: true,
    },
  };
  const filters = [
    { _id: 'filter-1-id', author: 'filter-1-author', title: 'filter-1-title' },
    { _id: 'filter-2-id', author: 'filter-2-author', title: 'filter-2-title' },
    { _id: 'filter-3-id', author: 'filter-3-author', title: 'filter-3-title' },
    { _id: 'filter-4-id', author: 'filter-4-author', title: 'filter-4-title' },
  ];

  const mockGetUserPreferenceRequest = (response = { content: {}, widget: widgetId, filters: [] }) =>
    axiosMockAdapter.onGet(`${API_ROUTES.userPreferences}/${widgetId}`).reply(200, response);

  const factory = generateShallowRenderer(FiltersList, {
    stubs,
    store,
    mocks: { $popups, $modals },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
        $modal: modal,
      },
    },
  });

  test('Filters fetched after mount', async () => {
    mockGetUserPreferenceRequest();

    factory({
      propsData: { modal },
    });

    await flushPromises();

    expect(fetchUserPreference).toBeDispatchedWith({ id: widgetId });
  });

  test('Filter created after trigger filters list', async () => {
    mockGetUserPreferenceRequest();
    createWidgetFilter.mockResolvedValueOnce({});

    const wrapper = factory({
      propsData: { modal },
    });

    await flushPromises();

    await selectFiltersList(wrapper).triggerCustomEvent('add');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createFilter,
      config: expect.objectContaining({
        action: expect.any(Function),
        corporate: true,
        withTitle: true,
        title: 'Create filter',
      }),
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const newFilter = { title: Faker.datatype.string() };

    await config.action(newFilter);

    expect(createWidgetFilter).toBeDispatchedWith({
      data: {
        ...newFilter,
        widget: widgetId,
        is_user_preference: true,
      },
    });
  });

  test('Filter edited after trigger filters list', async () => {
    mockGetUserPreferenceRequest();

    const wrapper = factory({
      propsData: { modal },
    });

    await flushPromises();

    const editingFilter = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      author: Faker.datatype.string(),
    };

    await selectFiltersList(wrapper).triggerCustomEvent('edit', editingFilter);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createFilter,
      config: expect.objectContaining({
        action: expect.any(Function),
        corporate: true,
        filter: editingFilter,
        withTitle: true,
        title: 'Edit filter',
      }),
    });

    const [{ config }] = $modals.show.mock.calls[0];
    const newFilter = { title: Faker.datatype.string() };

    updateWidgetFilter.mockResolvedValueOnce({});

    await config.action(newFilter);

    expect(updateWidgetFilter).toBeDispatchedWith({
      id: editingFilter._id,
      data: {
        ...newFilter,
        widget: widgetId,
      },
    });
  });

  test('Filter removed after trigger filters list', async () => {
    mockGetUserPreferenceRequest();

    const wrapper = factory({
      propsData: { modal },
    });

    await flushPromises();

    const removingFilter = {
      _id: Faker.datatype.string(),
    };

    await selectFiltersList(wrapper).triggerCustomEvent('delete', removingFilter);

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [{ config }] = $modals.show.mock.calls[0];

    removeWidgetFilter.mockResolvedValueOnce({});

    await config.action();

    expect(removeWidgetFilter).toBeDispatchedWith({ id: removingFilter._id });
  });

  test('Filters priority changed after trigger filters list', async () => {
    getUserPreferenceByWidgetId.mockReturnValue(() => ({
      content: {},
      widget: widgetId,
      filters,
    }));

    mockGetUserPreferenceRequest({
      content: {},
      widget: widgetId,
      filters,
    });

    const wrapper = factory({
      propsData: { modal },
    });

    await flushPromises();

    const newSortedFilters = [filters[0], filters[2], filters[1], filters[3]];

    updateWidgetFiltersPositions.mockResolvedValueOnce({});

    await selectFiltersList(wrapper).triggerCustomEvent('update:filters', newSortedFilters);

    await flushPromises();

    expect(updateWidgetFiltersPositions).toBeDispatchedWith({
      data: newSortedFilters.map(({ _id: id }) => id),
    });
  });

  test('Filters priority didn\'t change after trigger filters list with error', async () => {
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    getUserPreferenceByWidgetId.mockReturnValue(() => ({
      content: {},
      widget: widgetId,
      filters,
    }));

    mockGetUserPreferenceRequest({
      content: {},
      widget: widgetId,
      filters,
    });

    const wrapper = factory({
      propsData: { modal },
    });

    await flushPromises();

    wrapper.vm.setFilters([...filters]);

    const newSortedFilters = [filters[0], filters[2], filters[1], filters[3]];

    updateWidgetFiltersPositions.mockRejectedValueOnce(new Error('fail'));

    await selectFiltersList(wrapper).triggerCustomEvent('update:filters', newSortedFilters);

    await flushPromises();

    expect(updateWidgetFiltersPositions).toBeDispatchedWith({
      data: newSortedFilters.map(({ _id: id }) => id),
    });
    expect($popups.error).toHaveBeenCalledWith({ text: 'Something went wrong...' });
    expect(wrapper.vm.filters).toEqual(filters);

    consoleErrorSpy.mockRestore();
  });
});
