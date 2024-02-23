import Faker from 'faker';
import AxiosMockAdapter from 'axios-mock-adapter';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createModalWrapperStub } from '@unit/stubs/modal';

import { API_ROUTES } from '@/config';
import { MODALS } from '@/constants';

import ClickOutside from '@/services/click-outside';
import request from '@/services/request';

import store from '@/store';

import FiltersList from '@/components/modals/common/filters-list.vue';

const stubs = {
  'modal-wrapper': createModalWrapperStub('modal-wrapper'),
  'c-progress-overlay': true,
  'filters-list-component': true,
};

const selectFiltersList = wrapper => wrapper.find('filters-list-component-stub');

describe('filters-list', () => {
  mockDateNow(1386435600000);
  const $modals = mockModals();
  const $popups = mockPopups();

  const axiosMockAdapter = new AxiosMockAdapter(request);

  axiosMockAdapter
    .onAny()
    .passThrough();

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  const widgetId = 'widget-id';
  const modal = {
    config: {
      widgetId,
      filters: [],
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

  const mockGetUserPreferenceRequest = (response = {
    content: {},
    widget: widgetId,
    filters: [],
  }) => axiosMockAdapter
    .onGet(`${API_ROUTES.userPreferences}/${widgetId}`)
    .reply(200, response);
  const mockPostWidgetFiltersRequest = () => axiosMockAdapter
    .onPost(API_ROUTES.widget.filters)
    .reply(200);
  const mockPutWidgetFiltersRequest = id => axiosMockAdapter
    .onPut(`${API_ROUTES.widget.filters}/${id}`)
    .reply(200);
  const mockDeleteWidgetFiltersRequest = id => axiosMockAdapter
    .onDelete(`${API_ROUTES.widget.filters}/${id}`)
    .reply(200);
  const mockPutWidgetFilterPositionsRequest = () => axiosMockAdapter
    .onPut(API_ROUTES.widget.filterPositions)
    .reply(200);
  const mockPutWidgetFilterPositionsRequestWithError = () => axiosMockAdapter
    .onPut(API_ROUTES.widget.filterPositions)
    .reply(400);

  const factory = generateShallowRenderer(FiltersList, {
    stubs,
    store,
    mocks: { $popups, $modals },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(FiltersList, {
    stubs,
    store,
    mocks: { $popups, $modals },
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  test('Filters fetched after mount', async () => {
    mockGetUserPreferenceRequest();

    factory({
      propsData: {
        modal,
      },
    });

    await flushPromises();

    const [getRequestData] = axiosMockAdapter.history.get;

    expect(getRequestData).toBeTruthy();
  });

  test('Filter created after trigger filters list', async () => {
    mockGetUserPreferenceRequest();
    mockPostWidgetFiltersRequest();

    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    selectFiltersList(wrapper).triggerCustomEvent('add');

    expect($modals.show).toBeCalledWith({
      name: MODALS.createFilter,
      config: {
        action: expect.any(Function),
        corporate: true,
        withTitle: true,
        title: 'Create filter',
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const newFilter = { title: Faker.datatype.string() };

    await config.action(newFilter);

    const [filterPostRequest] = axiosMockAdapter.history.post;

    expect(JSON.parse(filterPostRequest.data)).toEqual({
      ...newFilter,
      widget: widgetId,
      is_user_preference: true,
    });
  });

  test('Filter edited after trigger filters list', async () => {
    mockGetUserPreferenceRequest();

    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    const editingFilter = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      author: Faker.datatype.string(),
    };

    selectFiltersList(wrapper).triggerCustomEvent('edit', editingFilter);

    expect($modals.show).toBeCalledWith({
      name: MODALS.createFilter,
      config: {
        action: expect.any(Function),
        filter: editingFilter,
        corporate: true,
        withTitle: true,
        title: 'Edit filter',
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const newFilter = { title: Faker.datatype.string() };

    mockPutWidgetFiltersRequest(editingFilter._id);

    await config.action(newFilter);

    const [filterPutRequest] = axiosMockAdapter.history.put;

    expect(JSON.parse(filterPutRequest.data)).toEqual({
      ...newFilter,
      widget: widgetId,
    });
  });

  test('Filter removed after trigger filters list', async () => {
    mockGetUserPreferenceRequest();

    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    const removingFilter = {
      _id: Faker.datatype.string(),
    };

    selectFiltersList(wrapper).triggerCustomEvent('delete', removingFilter);

    expect($modals.show).toBeCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    mockDeleteWidgetFiltersRequest(removingFilter._id);

    await config.action();

    const [filterDeleteRequest] = axiosMockAdapter.history.delete;

    expect(filterDeleteRequest).toBeTruthy();
  });

  test('Filters priority changed after trigger filters list', async () => {
    mockGetUserPreferenceRequest({
      filters,
      widget: widgetId,
    });

    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    const newSortedFilters = [
      filters[0],
      filters[2],
      filters[1],
      filters[3],
    ];

    mockPutWidgetFilterPositionsRequest();

    selectFiltersList(wrapper).triggerCustomEvent('input', newSortedFilters);

    await flushPromises();

    const [filterPutPositionsRequest] = axiosMockAdapter.history.put;

    expect(JSON.parse(filterPutPositionsRequest.data)).toEqual(
      newSortedFilters.map(({ _id: id }) => id),
    );
  });

  test('Filters priority didn\'t change after trigger filters list with error', async () => {
    const consoleErrorSpy = jest
      .spyOn(console, 'error')
      .mockImplementation(() => {});

    mockGetUserPreferenceRequest({
      filters,
      widget: widgetId,
    });

    const wrapper = factory({
      propsData: {
        modal,
      },
    });

    const newSortedFilters = [
      filters[0],
      filters[2],
      filters[1],
      filters[3],
    ];

    mockPutWidgetFilterPositionsRequestWithError();

    selectFiltersList(wrapper).triggerCustomEvent('input', newSortedFilters);

    await flushPromises();

    const [filterPutPositionsRequest] = axiosMockAdapter.history.put;

    expect(JSON.parse(filterPutPositionsRequest.data)).toEqual(
      newSortedFilters.map(({ _id: id }) => id),
    );

    expect($popups.error).toBeCalledWith({ text: 'Something went wrong...' });

    consoleErrorSpy.mockClear();
  });

  test('Renders `filters-list` with empty modal', async () => {
    mockGetUserPreferenceRequest();

    const wrapper = snapshotFactory({
      store,
      propsData: {
        modal: {
          config: {
            widgetId,
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `filters-list` with all parameters', async () => {
    mockGetUserPreferenceRequest({
      content: {},
      widget: widgetId,
      filters,
    });
    const wrapper = snapshotFactory({
      propsData: {
        modal,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
