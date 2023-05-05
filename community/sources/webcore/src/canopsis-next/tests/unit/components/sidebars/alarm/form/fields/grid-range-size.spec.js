import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createNumberInputStub } from '@unit/stubs/input';
import GridRangeSize from '@/components/sidebars/alarm/form/fields/grid-range-size.vue';

const stubs = {
  'v-range-slider': createNumberInputStub('v-range-slider'),
};

const factory = generateShallowRenderer(GridRangeSize, { stubs,
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

const snapshotFactory = generateRenderer(GridRangeSize, { stubs,
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

const selectRangeSliderField = wrapper => wrapper.find('input.v-range-slider');

describe('grid-range-size', () => {
  it('Value changed after trigger range slider field', () => {
    const wrapper = factory();

    const rangeSliderField = selectRangeSliderField(wrapper);

    const newValue = 3;

    rangeSliderField.setValue(newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(newValue);
  });

  it('Renders `grid-range-size` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `grid-range-size` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [2, 10],
        min: 2,
        max: 10,
        step: 2,
        title: 'Custom title',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
