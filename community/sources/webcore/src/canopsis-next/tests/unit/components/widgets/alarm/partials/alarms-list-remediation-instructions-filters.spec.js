import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { mockModals } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { MODALS } from '@/constants';

import AlarmsListRemediationInstructionsFilters from '@/components/widgets/alarm/partials/alarms-list-remediation-instructions-filters.vue';

const stubs = {
  'remediation-instructions-filters-list': true,
  'v-btn': createButtonStub('v-btn'),
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

const snapshotStubs = {
  'remediation-instructions-filters-list': true,
};

const updateFilters = jest.fn();
const updateLockedFilters = jest.fn();

const selectLockedRemediationInstructionsFiltersList = wrapper => wrapper.findAll('remediation-instructions-filters-list-stub').at(0);
const selectRemediationInstructionsFiltersList = wrapper => wrapper.findAll('remediation-instructions-filters-list-stub').at(1);
const selectAddFilterButton = wrapper => wrapper.find('button.v-btn');

describe('alarms-list-remediation-instructions-filters', () => {
  const $modals = mockModals();

  const lockedFilters = [
    {
      title: 'Locked filter 1',
      _id: 'ID locked filter 1',
    },
    {
      title: 'Locked filter 2',
      _id: 'ID locked filter 2',
    },
  ];
  const filters = [
    {
      title: 'Filter 1',
      _id: 'ID filter 1',
    },
    {
      title: 'Filter 2',
      _id: 'ID filter 2',
    },
  ];

  const factory = generateShallowRenderer(
    AlarmsListRemediationInstructionsFilters,
    {
      stubs,
      listeners: {
        'update:filters': updateFilters,
        'update:lockedFilters': updateLockedFilters,
      },
    },
  );
  const snapshotFactory = generateRenderer(
    AlarmsListRemediationInstructionsFilters,
    {
      stubs: snapshotStubs,
      listeners: {
        'update:filters': updateFilters,
        'update:lockedFilters': updateLockedFilters,
      },
    },
  );

  afterEach(() => {
    updateFilters.mockReset();
    updateLockedFilters.mockReset();
  });

  it('Locked filters updated after trigger remediation instructions filters list', () => {
    const wrapper = factory({
      propsData: {
        lockedFilters: [],
      },
    });

    const lockedRemediationInstructionsFiltersList = selectLockedRemediationInstructionsFiltersList(wrapper);

    lockedRemediationInstructionsFiltersList.vm.$emit('input', lockedFilters);

    expect(updateLockedFilters).toHaveBeenCalledTimes(1);
    expect(updateLockedFilters).toHaveBeenCalledWith(lockedFilters);
  });

  it('Filters updated after trigger remediation instructions filters list', () => {
    const wrapper = factory({
      propsData: {
        filters: [],
      },
    });

    const remediationInstructionsFiltersList = selectRemediationInstructionsFiltersList(wrapper);

    remediationInstructionsFiltersList.vm.$emit('input', filters);

    expect(updateFilters).toHaveBeenCalledTimes(1);
    expect(updateFilters).toHaveBeenCalledWith(filters);
  });

  it('Create remediation instruction filter modal opened after trigger a button', () => {
    const wrapper = factory({
      propsData: {
        addable: true,
        filters,
      },
      mocks: {
        $modals,
      },
    });

    const addButton = selectAddFilterButton(wrapper);

    addButton.trigger('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createRemediationInstructionsFilter,
        config: {
          filters,
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = {
      title: 'New filter',
      filter: {},
    };

    modalArguments.config.action(actionValue);

    const updateFiltersEvents = wrapper.emitted('update:filters');

    expect(updateFiltersEvents).toHaveLength(1);

    const [eventData] = updateFiltersEvents[0];
    expect(eventData).toEqual([
      ...filters,
      {
        ...actionValue,
        _id: expect.any(String),
      },
    ]);
  });

  it('Renders `alarms-list-remediation-instructions-filters` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-remediation-instructions-filters` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        lockedFilters,
        addable: true,
        editable: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-remediation-instructions-filters` with locked filters props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        lockedFilters,
        addable: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-remediation-instructions-filters` with access, but without filters filters props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        addable: true,
        editable: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
