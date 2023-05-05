import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import MermaidThemeField from '@/components/other/map/form/fields/mermaid-theme-field.vue';
import { MERMAID_THEMES } from '@/constants';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = generateShallowRenderer(MermaidThemeField, { stubs,
});

const snapshotFactory = generateRenderer(MermaidThemeField, {

});

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('mermaid-theme-field', () => {
  test('Value changed after trigger select', () => {
    const wrapper = factory({
      propsData: {
        value: MERMAID_THEMES.canopsis,
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', MERMAID_THEMES.default);

    expect(wrapper).toEmit('input', MERMAID_THEMES.default);
  });

  test('Renders `mermaid-theme-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: MERMAID_THEMES.canopsis,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
