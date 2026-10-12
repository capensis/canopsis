import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { MERMAID_THEMES } from '@/constants';

import MermaidThemeField from '@/components/other/map/form/fields/mermaid-theme-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('select.v-select');

describe('mermaid-theme-field', () => {
  const factory = generateShallowRenderer(MermaidThemeField, { stubs });

  test('Value changed after trigger select', () => {
    const wrapper = factory({
      propsData: {
        value: MERMAID_THEMES.canopsis,
      },
    });

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', MERMAID_THEMES.default);

    expect(wrapper).toEmitInput(MERMAID_THEMES.default);
  });
});
