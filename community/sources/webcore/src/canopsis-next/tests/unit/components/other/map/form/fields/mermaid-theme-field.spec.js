import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { MERMAID_THEMES } from '@/constants';

import MermaidThemeField from '@/components/other/map/form/fields/mermaid-theme-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('mermaid-theme-field', () => {
  const factory = generateShallowRenderer(MermaidThemeField, { stubs });
  const snapshotFactory = generateRenderer(MermaidThemeField);

  test('Value changed after trigger select', () => {
    const wrapper = factory({
      propsData: {
        value: MERMAID_THEMES.canopsis,
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', MERMAID_THEMES.default);

    expect(wrapper).toEmit('input', MERMAID_THEMES.default);
  });

  test('Renders `mermaid-theme-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: MERMAID_THEMES.canopsis,
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `mermaid-theme-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: MERMAID_THEMES.dark,
        label: 'Custom label',
        name: 'custom_name',
      },
    });

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
