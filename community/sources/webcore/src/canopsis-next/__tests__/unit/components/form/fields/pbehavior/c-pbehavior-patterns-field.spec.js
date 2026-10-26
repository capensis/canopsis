import { generateShallowRenderer } from '@unit/utils/vue';

import { PBEHAVIOR_PATTERN_FIELDS } from '@/constants';

import CPbehaviorPatternsField from '@/components/forms/fields/pbehavior/c-pbehavior-patterns-field.vue';

const patternEditorFieldStub = {
  name: 'PatternEditorField',
  props: ['patterns', 'attributes'],
  template: '<div />',
};

const stubs = {
  'pattern-editor-field': patternEditorFieldStub,
};

describe('c-pbehavior-patterns-field', () => {
  const factory = generateShallowRenderer(CPbehaviorPatternsField, { stubs });

  test('Uses lazy fields for PBehavior attributes', () => {
    const wrapper = factory({
      propsData: {
        patterns: {},
      },
    });
    const attributes = wrapper.findComponent(patternEditorFieldStub).props('attributes');
    const attributesByValue = Object.fromEntries(attributes.map(attribute => [attribute.value, attribute]));

    expect(attributesByValue[PBEHAVIOR_PATTERN_FIELDS.id].options.valueField).toEqual({
      is: 'c-pbehavior-field',
      model: 'modelValue',
      props: {
        itemTitle: 'name',
        itemValue: '_id',
        ellipsis: true,
      },
    });
    expect(attributesByValue[PBEHAVIOR_PATTERN_FIELDS.reason].options.valueField).toEqual({
      is: 'c-pbehavior-reason-field',
      model: 'modelValue',
      props: {
        ellipsis: true,
      },
    });
    expect(attributesByValue[PBEHAVIOR_PATTERN_FIELDS.type].options.valueField).toEqual({
      is: 'c-pbehavior-type-field',
      model: 'modelValue',
      props: {
        ellipsis: true,
      },
    });
  });
});
