import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { COLORS } from '@/config';

import TagGeneralForm from '@/components/other/tag/form/tag-general-form.vue';

const stubs = {
  'c-form-block': true,
  'c-form-block-row': true,
  'c-name-field': true,
  'c-color-picker-field': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectColorField = wrapper => wrapper.find('c-color-picker-field-stub');

describe('tag-general-form', () => {
  const form = {
    value: Faker.datatype.string(),
    color: Faker.internet.color(),
  };

  const factory = generateShallowRenderer(TagGeneralForm, { stubs });
  const snapshotFactory = generateRenderer(TagGeneralForm, { stubs });

  test('Value changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newValue = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ ...form, value: newValue });
  });

  test('Color changed after trigger color field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newColor = Faker.internet.color();

    selectColorField(wrapper).triggerCustomEvent('input', newColor);

    expect(wrapper).toEmitInput({ ...form, color: newColor });
  });

  test('Renders `tag-general-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `tag-general-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          value: 'tag-value',
          color: COLORS.secondary,
        },
        isImported: true,
        isNew: false,
        maxTagNameLength: 11,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
