import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

export const generateEntityPatternsTests = (Component, name, customProps = {}) => {
  const stubs = {
    'c-patterns-field': true,
  };

  const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

  describe(name, () => {
    const factory = generateShallowRenderer(Component, { stubs });
    const snapshotFactory = generateRenderer(Component, { stubs });

    test('Patterns changed after trigger patterns field', () => {
      const wrapper = factory();

      const patternsField = selectPatternsField(wrapper);

      const newPatterns = {
        alarm_pattern: {},
        pbehavior_pattern: {},
        entity_pattern: {},
      };

      patternsField.triggerCustomEvent('input', newPatterns);

      expect(wrapper).toEmit('input', newPatterns);
    });

    test(`Renders \`${name}\` with default props`, () => {
      const wrapper = snapshotFactory();

      expect(wrapper).toMatchSnapshot();
    });

    test(`Renders \`${name}\` with custom props`, () => {
      const wrapper = snapshotFactory({
        propsData: {
          form: {
            alarm_pattern: {},
            entity_pattern: {},
          },
          ...customProps,
        },
      });

      expect(wrapper).toMatchSnapshot();
    });
  });
};
