import { generateShallowRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import RootCauseSettings from '@/components/sidebars/form/fields/root-cause-settings.vue';

const stubs = {
  'widget-settings-item': true,
  'v-expand-transition': {
    template: '<div><slot /></div>',
  },
  'c-enabled-field': true,
  'c-color-indicator-field': true,
};

describe('field-color-indicator', () => {
  const factory = generateShallowRenderer(RootCauseSettings, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const modelValue = {
      showRootCauseByStateClick: true,
      rootCauseColorIndicator: COLOR_INDICATOR_TYPES.state,
    };
    const wrapper = factory({
      propsData: {
        value: modelValue,
      },
    });

    wrapper.vm.updateField('rootCauseColorIndicator', COLOR_INDICATOR_TYPES.impactState);

    expect(wrapper).toEmitInput({
      ...modelValue,
      rootCauseColorIndicator: COLOR_INDICATOR_TYPES.impactState,
    });
  });
});
