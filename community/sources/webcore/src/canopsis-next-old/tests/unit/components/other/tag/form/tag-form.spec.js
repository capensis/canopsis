import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';

import { COLORS } from '@/config';
import { IDLE_RULE_TYPES } from '@/constants';

import TagForm from '@/components/other/tag/form/tag-form.vue';

const stubs = {
  'tag-general-form': true,
  'tag-patterns-form': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
};

const selectTagGeneralForm = wrapper => wrapper.find('tag-general-form-stub');
const selectTagPatternsForm = wrapper => wrapper.find('tag-patterns-form-stub');

describe('tag-form', () => {
  const factory = generateShallowRenderer(TagForm, { stubs });
  const snapshotFactory = generateRenderer(TagForm, { stubs });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form: {
          value: '',
        },
      },
    });

    expect(selectTagGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Value changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          value: '',
        },
      },
    });

    const newValue = Faker.datatype.string();

    selectTagGeneralForm(wrapper).triggerCustomEvent('input', { value: newValue });

    expect(wrapper).toEmitInput({ value: newValue });
  });

  test('Color changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          color: Faker.internet.color(),
        },
      },
    });

    const newValue = Faker.internet.color();

    selectTagGeneralForm(wrapper).triggerCustomEvent('input', { color: newValue });

    expect(wrapper).toEmitInput({ color: newValue });
  });

  test('Tag patterns changed after trigger patterns form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          value: 'Value',
          patterns: {},
        },
      },
    });

    const newPatterns = {
      alarm_pattern: {},
      entity_pattern: {},
    };

    selectTagPatternsForm(wrapper).triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmitInput({
      value: 'Value',
      patterns: newPatterns,
    });
  });

  test('Renders `tag-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
