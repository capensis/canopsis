import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';
import { mockModals } from '@unit/utils/mock-hooks';
import { MODALS } from '@/constants';

import RemediationInstructionsFilters from '@/components/side-bars/settings/fields/common/remediation-instructions-filters.vue';

const localVue = createVueInstance();

const stubs = {
  'remediation-instructions-filters-list': true,
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'remediation-instructions-filters-list': true,
};

const factory = (options = {}) => shallowMount(RemediationInstructionsFilters, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(RemediationInstructionsFilters, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    provide: {
      list: {
        register: jest.fn(),
        unregister: jest.fn(),
      },
      listClick: jest.fn(),
    },
  },

  ...options,
});

const selectRemediationInstructionsFiltersListField = wrapper => wrapper.find('remediation-instructions-filters-list-stub');
const selectAddButton = wrapper => wrapper.find('button.v-btn');

describe('remediation-instructions-filters', () => {
  const $modals = mockModals();
  const filters = [{
    with: true,
    all: true,
    auto: true,
    manual: true,
    instructions: [],
    _id: 'id',
  }];

  it('Instruction filters changed after trigger separator select field', () => {
    const wrapper = factory({
      propsData: {
        filters: [],
      },
    });

    const remediationInstructionsFiltersListField = selectRemediationInstructionsFiltersListField(wrapper);

    remediationInstructionsFiltersListField.vm.$emit('input', filters);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(filters);
  });

  it('Instruction filter added after separator add button', () => {
    const wrapper = factory({
      propsData: {
        filters,
        addable: true,
      },
      mocks: {
        $modals,
      },
    });

    const addButton = selectAddButton(wrapper);

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
      with: false,
      all: false,
      auto: false,
      manual: false,
      instructions: [],
    };

    modalArguments.config.action(actionValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual([
      ...filters,
      { ...actionValue, _id: expect.any(String) },
    ]);
  });

  it('Renders `remediation-instructions-filters` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `remediation-instructions-filters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        editable: true,
        addable: true,
      },
    });

    const menuContents = wrapper.findAll('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    menuContents.wrappers.forEach((menuContent) => {
      expect(menuContent.element).toMatchSnapshot();
    });
  });
});