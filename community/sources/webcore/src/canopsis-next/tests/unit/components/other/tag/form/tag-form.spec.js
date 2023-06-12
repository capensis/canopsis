import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { IDLE_RULE_TYPES } from '@/constants';

import TagForm from '@/components/other/tag/form/tag-form.vue';
import { COLORS } from '@/config';

const stubs = {
  'c-name-field': true,
  'c-color-picker-field': true,
  'tag-patterns-form': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectColorPickerField = wrapper => wrapper.find('c-color-picker-field-stub');
const selectTagPatternsForm = wrapper => wrapper.find('tag-patterns-form-stub');

describe('tag-form', () => {
  const factory = generateShallowRenderer(TagForm, { stubs });
  const snapshotFactory = generateRenderer(TagForm, { stubs });

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          name: '',
        },
      },
    });

    const newValue = Faker.datatype.string();

    selectNameField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', { name: newValue });
  });

  test('Color changed after trigger color picker field', () => {
    const wrapper = factory({
      propsData: {
        form: {
          color: Faker.internet.color(),
        },
      },
    });

    const newValue = Faker.internet.color();

    selectColorPickerField(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', { color: newValue });
  });

  test('Tag patterns changed after trigger patterns form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          name: 'Name',
          patterns: {},
        },
      },
    });

    const newPatterns = {
      alarm_pattern: {},
      entity_pattern: {},
    };

    selectTagPatternsForm(wrapper).vm.$emit('input', newPatterns);

    expect(wrapper).toEmit('input', {
      name: 'Name',
      patterns: newPatterns,
    });
  });

  test('Renders `tag-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `tag-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          type: IDLE_RULE_TYPES.entity,
          color: COLORS.secondary,
          patterns: {},
        },
        isImported: true,
        maxTagNameLength: 11,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
