import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createNumberInputStub } from '@unit/stubs/input';

import FieldSlider from '@/components/sidebars/form/fields/slider.vue';

const stubs = {
  'widget-settings-item': true,
  'v-slider': createNumberInputStub('v-slider'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectSliderField = wrapper => wrapper.find('.v-slider');

describe('field-slider', () => {
  const factory = generateShallowRenderer(FieldSlider, { stubs });
  const snapshotFactory = generateRenderer(FieldSlider, { stubs: snapshotStubs });

  test('Value changed after trigger number field', () => {
    const value = Faker.datatype.number();
    const wrapper = factory({
      propsData: {
        title: '',
        value,
      },
    });

    const newValue = Faker.datatype.number({ min: value });

    selectSliderField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `field-slider` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 3,
        title: 'Custom required title',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `field-slider` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 22,
        title: 'Custom title',
        min: 20,
        max: 25,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
