import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';
import { mockModals } from '@unit/utils/mock-hooks';
import { MODALS } from '@/constants';

import RemediationInstructionsFiltersItem from '@/components/other/remediation/instructions-filter/partials/remediation-instructions-filters-item.vue';

const stubs = {
  'v-chip': createButtonStub('v-chip'),
};

const selectChip = wrapper => wrapper.find('button.v-chip');

describe('remediation-instructions-filters-item', () => {
  const $modals = mockModals();
  const lockedFilter = {
    with: true,
    all: false,
    auto: true,
    manual: false,
    locked: true,
    disabled: false,
    instructions: [{ name: 'instruction-1' }],
    _id: 'id1',
  };
  const unLockedFilter = {
    with: false,
    all: true,
    auto: false,
    manual: false,
    locked: false,
    disabled: true,
    instructions: [],
    _id: 'id2',
  };
  const filters = [lockedFilter, unLockedFilter];

  const factory = generateShallowRenderer(RemediationInstructionsFiltersItem, { stubs });
  const snapshotFactory = generateRenderer(RemediationInstructionsFiltersItem, {
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

  it('Locked filter disabled after trigger input event on the chip', () => {
    const wrapper = factory({
      propsData: {
        filter: lockedFilter,
      },
    });

    const chip = selectChip(wrapper);

    chip.vm.$emit('input');

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...lockedFilter,
      disabled: !lockedFilter.disabled,
    });
  });

  it('Unlocked filter removed after trigger input event on the chip', () => {
    const wrapper = factory({
      propsData: {
        filter: unLockedFilter,
        editable: true,
      },
    });

    const chip = selectChip(wrapper);

    chip.vm.$emit('input');

    const inputEvents = wrapper.emitted('remove');

    expect(inputEvents).toHaveLength(1);
  });

  it('Edit instruction filter modal opened after trigger click event on the chip', () => {
    const wrapper = factory({
      propsData: {
        filters,
        filter: unLockedFilter,
        editable: true,
      },
      mocks: {
        $modals,
      },
    });

    const chip = selectChip(wrapper);

    chip.vm.$emit('click');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createRemediationInstructionsFilter,
        config: {
          filter: unLockedFilter,
          filters: [lockedFilter],
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
    expect(eventData).toEqual(actionValue);
  });

  it('Renders `remediation-instructions-filters-item` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `remediation-instructions-filters-item` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        filter: {
          with: true,
          all: false,
          auto: false,
          manual: true,
          locked: true,
          disabled: true,
          instructions: [],
          _id: 'id3',
        },
        editable: true,
        closable: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `remediation-instructions-filters-item` with filter in progress', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        filter: {
          with: true,
          all: false,
          auto: false,
          manual: true,
          locked: true,
          disabled: true,
          instructions: [],
          running: true,
          _id: 'id3',
        },
        editable: true,
        closable: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `remediation-instructions-filters-item` with filter instruction in progress', () => {
    const wrapper = snapshotFactory({
      propsData: {
        filters,
        filter: {
          with: true,
          all: false,
          auto: false,
          manual: true,
          locked: true,
          disabled: true,
          instructions: [],
          running: true,
          _id: 'id3',
        },
        editable: true,
        closable: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
