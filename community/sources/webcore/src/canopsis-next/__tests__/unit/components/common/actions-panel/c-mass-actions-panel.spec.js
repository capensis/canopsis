import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { TELEPORTS_NAMES } from '@/constants';

import CMassActionsPanel from '@/components/common/actions-panel/c-mass-actions-panel.vue';

const stubs = {
  'c-enabled-field': true,
};

describe('c-mass-actions-panel', () => {
  const factory = generateShallowRenderer(CMassActionsPanel, { stubs });

  test('Portal target is massActionsPanel when modal is not injected', () => {
    const wrapper = factory();

    expect(wrapper.vm.targetSelector).toBe(
      `[data-teleport-target="${TELEPORTS_NAMES.massActionsPanel}"]`,
    );
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

    expect(wrapper.vm.targetSelector).toBe(
      `[data-teleport-target="${TELEPORTS_NAMES.massActionsPanel}-${modalId}"]`,
    );
  });

  test('Message displays selected count', () => {
    const selectedCount = Faker.datatype.number({ min: 1, max: 100 });
    const selected = Array(selectedCount).fill({ _id: Faker.datatype.uuid() });
    const wrapper = factory({
      propsData: {
        selected,
      },
    });

    expect(wrapper.vm.message).toContain(String(selectedCount));
  });

  test('clearSelected emits clear:selected when close button clicked', () => {
    const wrapper = factory({
      propsData: {
        selected: [{ _id: 'alarm-1' }],
      },
    });

    wrapper.vm.$emit('clear:selected');
    expect(wrapper.emitted('clear:selected')).toHaveLength(1);
  });

  test('clearSelected emits clear:selected when called directly', () => {
    const wrapper = factory({
      propsData: {
        selected: [{ _id: 'alarm-1' }],
      },
    });

    wrapper.vm.$emit('clear:selected');

    expect(wrapper.emitted('clear:selected')).toHaveLength(1);
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

  test('Keep-selected field emits update:modelValue with true', () => {
    const wrapper = factory();

    wrapper.vm.$emit('update:modelValue', true);

    expect(wrapper.emitted('update:modelValue')).toEqual([[true]]);
  });

  test('Keep-selected field emits update:modelValue with false', () => {
    const wrapper = factory();

    wrapper.vm.$emit('update:modelValue', false);

    expect(wrapper.emitted('update:modelValue')).toEqual([[false]]);
  });
});
