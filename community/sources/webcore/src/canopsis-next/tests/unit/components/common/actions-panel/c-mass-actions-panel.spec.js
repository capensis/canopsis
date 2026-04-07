import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';

import { PORTALS_NAMES } from '@/constants';

import CMassActionsPanel from '@/components/common/actions-panel/c-mass-actions-panel.vue';

const stubs = {
  portal: true,
  'c-enabled-field': true,
  'c-action-btn': createButtonStub('c-action-btn'),
};

const snapshotStubs = {
  portal: {
    props: ['to'],
    template: `
      <div class="portal">
        Portal to: {{ to }}
        <slot />
      </div>
    `,
  },
  'c-enabled-field': true,
  'c-action-btn': createButtonStub('c-action-btn'),
};

const selectPortal = wrapper => wrapper.find('portal-stub');
const selectClearButton = wrapper => wrapper.find('button.c-action-btn');

describe('c-mass-actions-panel', () => {
  const factory = generateShallowRenderer(CMassActionsPanel, { stubs });
  const snapshotFactory = generateRenderer(CMassActionsPanel, { stubs: snapshotStubs });

  test('Portal target is massActionsPanel when modal is not injected', () => {
    const wrapper = factory();

    const portal = selectPortal(wrapper);

    expect(portal.attributes('to')).toBe(PORTALS_NAMES.massActionsPanel);
  });

  test('Portal target includes modal id when modal is injected', () => {
    const modalId = Faker.datatype.uuid();
    const wrapper = factory({
      parentComponent: {
        provide: {
          $modal: {
            id: modalId,
          },
        },
      },
    });

    const portal = selectPortal(wrapper);

    expect(portal.attributes('to')).toBe(`${PORTALS_NAMES.massActionsPanel}-${modalId}`);
  });

  test('Message displays selected count', () => {
    const selectedCount = Faker.datatype.number({ min: 1, max: 100 });
    const selected = Array(selectedCount).fill({ _id: Faker.datatype.uuid() });
    const wrapper = factory({
      propsData: {
        selected,
      },
    });

    expect(wrapper.text()).toContain(String(selectedCount));
  });

  test('clearSelected emits clear:selected when close button clicked', () => {
    const wrapper = factory({
      propsData: {
        selected: [{ _id: 'alarm-1' }],
      },
    });

    const clearButton = selectClearButton(wrapper);

    clearButton.trigger('click');

    expect(wrapper).toHaveBeenEmit('clear:selected');
  });

  test('clearSelected emits clear:selected when called directly', () => {
    const wrapper = factory({
      propsData: {
        selected: [{ _id: 'alarm-1' }],
      },
    });

    wrapper.vm.clearSelected();

    expect(wrapper).toHaveBeenEmit('clear:selected');
  });

  test('Style computed from grid parameters', () => {
    const x = 2;
    const y = 1;
    const w = 6;
    const wrapper = factory({
      propsData: {
        x,
        y,
        w,
      },
    });

    expect(wrapper.vm.style).toEqual({
      gridColumnStart: x + 1,
      gridColumnEnd: x + 1 + w,
      gridRowStart: y + 1,
      gridRowEnd: y + 2,
    });
  });

  test('updateKeepSelected emits update:keep-selected with true', () => {
    const wrapper = factory();

    wrapper.vm.updateKeepSelected(true);

    expect(wrapper).toEmit('update:keep-selected', true);
  });

  test('updateKeepSelected emits update:keep-selected with false', () => {
    const wrapper = factory();

    wrapper.vm.updateKeepSelected(false);

    expect(wrapper).toEmit('update:keep-selected', false);
  });

  test('Renders `c-mass-actions-panel` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: [{ _id: 'alarm-1' }],
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper.text()).toContain(PORTALS_NAMES.massActionsPanel);
  });

  test('Renders `c-mass-actions-panel` with modal id', () => {
    const modalId = 'alarms-list-modal';
    const wrapper = snapshotFactory({
      propsData: {
        selected: [{ _id: 'alarm-1' }, { _id: 'alarm-2' }],
      },
      parentComponent: {
        provide: {
          $modal: {
            id: modalId,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper.text()).toContain(`${PORTALS_NAMES.massActionsPanel}-${modalId}`);
  });

  test('Renders `c-mass-actions-panel` with actions slot', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: [{ _id: 'alarm-1' }],
      },
      slots: {
        actions: '<div class="custom-actions">Custom actions</div>',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper.find('.custom-actions').text()).toBe('Custom actions');
  });
});
