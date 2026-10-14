import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import GridRangeSize from '@/components/sidebars/form/fields/grid-range-size.vue';

const stubs = {
  'v-range-slider': createNumberInputStub('v-range-slider'),
};

const selectRangeSliderField = wrapper => wrapper.find('input.v-range-slider');

describe('grid-range-size', () => {
  const factory = generateShallowRenderer(GridRangeSize, {
    stubs,
    parentComponent: {
      provide: {
        list: {
          register: vi.fn(),
          unregister: vi.fn(),
          listClick: vi.fn(),
        },
      },
    },
  });

  it('Value changed after trigger range slider field', () => {
    const wrapper = factory();

    const newValue = 3;

    selectRangeSliderField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
