import { generateShallowRenderer, generateRenderer, flushPromises } from '@unit/utils/vue';
import { createNumberInputStub } from '@unit/stubs/input';

import GridRangeSize from '@/components/sidebars/form/fields/grid-range-size.vue';

const stubs = {
  'v-range-slider': createNumberInputStub('v-range-slider'),
};

const selectRangeSliderField = wrapper => wrapper.find('input.v-range-slider');
const selectListTitle = wrapper => wrapper.find('.v-list-item__title');

describe('grid-range-size', () => {
  const factory = generateShallowRenderer(GridRangeSize, { stubs,
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
          listClick: jest.fn(),
        },
      },
    },
  });
  const snapshotFactory = generateRenderer(GridRangeSize, {
    stubs,
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
          listClick: jest.fn(),
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

  it('Renders `grid-range-size` with default and required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `grid-range-size` with default and required props (opened)', async () => {
    const wrapper = snapshotFactory();

    selectListTitle(wrapper).trigger('click');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `grid-range-size` with custom props (opened)', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [2, 10],
        min: 2,
        max: 10,
        step: 2,
        title: 'Custom title',
      },
    });

    selectListTitle(wrapper).trigger('click');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
