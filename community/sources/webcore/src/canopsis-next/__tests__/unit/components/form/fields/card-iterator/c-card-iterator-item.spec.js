import { generateShallowRenderer } from '@unit/utils/vue';

import CCardIteratorItem from '@/components/forms/fields/card-iterator/c-card-iterator-item.vue';

const stubs = {
  'v-row': true,
  'v-col': true,
  'v-expand-transition': true,
  'c-draggable-step-number': {
    name: 'CDraggableStepNumber',
    props: ['disabled'],
    template: '<div><slot /></div>',
  },
  'c-action-btn': true,
};

const selectDraggableStepNumber = wrapper =>
  wrapper.findComponent({ name: 'CDraggableStepNumber' });

describe('c-card-iterator-item', () => {
  const factory = generateShallowRenderer(CCardIteratorItem, { stubs });

  test.each([false, true])('Passes disabled state %s to the drag handle', (disabled) => {
    const wrapper = factory({ props: { disabled } });

    expect(selectDraggableStepNumber(wrapper).props('disabled')).toBe(disabled);
  });
});
