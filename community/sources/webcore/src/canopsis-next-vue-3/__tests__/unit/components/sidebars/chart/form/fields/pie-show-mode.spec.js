import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import CPieShowModeField from '@/components/sidebars/chart/form/fields/pie-show-mode.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('pie-show-mode', () => {
  const factory = generateShallowRenderer(CPieShowModeField, { stubs });

  test('Value changed after trigger radio group field', () => {
    const wrapper = factory({
      propsData: {
        value: KPI_PIE_CHART_SHOW_MODS.numbers,
      },
    });

    selectRadioGroup(wrapper).setValue(KPI_PIE_CHART_SHOW_MODS.percent);

    expect(wrapper).toEmitInput(KPI_PIE_CHART_SHOW_MODS.percent);
  });
});
