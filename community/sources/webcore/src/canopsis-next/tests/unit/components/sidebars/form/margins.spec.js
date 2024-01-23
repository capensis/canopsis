import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import Margins from '@/components/sidebars/form/margins.vue';

const stubs = {
  'widget-settings-group': true,
  'field-slider': true,
};

const selectFieldSliderByIndex = (wrapper, index) => wrapper.findAll('field-slider-stub').at(index);
const selectFieldSliderTop = wrapper => selectFieldSliderByIndex(wrapper, 0);
const selectFieldSliderRight = wrapper => selectFieldSliderByIndex(wrapper, 1);
const selectFieldSliderBottom = wrapper => selectFieldSliderByIndex(wrapper, 2);
const selectFieldSliderLeft = wrapper => selectFieldSliderByIndex(wrapper, 3);

describe('margins', () => {
  const form = {
    top: Faker.datatype.number(),
    right: Faker.datatype.number(),
    bottom: Faker.datatype.number(),
    left: Faker.datatype.number(),
  };

  const factory = generateShallowRenderer(Margins, {

    stubs,
  });

  const snapshotFactory = generateRenderer(Margins, {

    stubs,
  });

  test('Top margin changed after trigger slider field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newValue = Faker.datatype.number();

    selectFieldSliderTop(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, top: newValue });
  });

  test('Right margin changed after trigger slider field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newValue = Faker.datatype.number();

    selectFieldSliderRight(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, right: newValue });
  });

  test('Bottom margin changed after trigger slider field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newValue = Faker.datatype.number();

    selectFieldSliderBottom(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, bottom: newValue });
  });

  test('Left margin changed after trigger slider field', () => {
    const wrapper = factory({
      propsData: { form },
    });

    const newValue = Faker.datatype.number();

    selectFieldSliderLeft(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, left: newValue });
  });

  test('Renders `margins` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          top: 7,
          right: 5,
          bottom: 3,
          left: 1,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `margins` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          top: 1,
          right: 3,
          bottom: 5,
          left: 7,
        },
        min: 1,
        max: 10,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
