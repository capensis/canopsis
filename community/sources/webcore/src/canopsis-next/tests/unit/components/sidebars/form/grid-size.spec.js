import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import FieldGridSize from '@/components/sidebars/form/fields/grid-size.vue';

const stubs = {
  'widget-settings-item': true,
  'c-column-size-field': true,
};

const selectColumnSizeField = wrapper => wrapper.find('c-column-size-field-stub');

describe('field-grid-size', () => {
  const factory = generateShallowRenderer(FieldGridSize, { stubs });
  const snapshotFactory = generateRenderer(FieldGridSize, { stubs });

  it('Column size changed after trigger field', () => {
    const wrapper = factory({
      propsData: {
        title: 'Title',
      },
    });

    const columnSizeField = selectColumnSizeField(wrapper);

    const newSize = Faker.datatype.number({ min: 1, max: 12 });

    columnSizeField.triggerCustomEvent('input', newSize);

    expect(wrapper).toEmitInput(newSize);
  });

  it('Renders `field-grid-size` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Custom title',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `field-grid-size` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 12,
        title: 'Custom title',
        mobile: true,
        tablet: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
