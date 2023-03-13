import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { KPI_PIE_CHART_SHOW_MODES } from '@/constants';

import CPieShowModeField from '@/components/sidebars/settings/fields/chart/pie-show-mode.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createInputStub('v-radio-group'),
};
const snapshotStubs = {
  'widget-settings-item': true,
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('pie-show-mode', () => {
  const factory = generateShallowRenderer(CPieShowModeField, {
    stubs,
  });
  const snapshotFactory = generateRenderer(CPieShowModeField, {
    stubs: snapshotStubs,
  });

  test('Value changed after trigger radio group field', () => {
    const wrapper = factory({
      propsData: {
        value: KPI_PIE_CHART_SHOW_MODES.numbers,
      },
    });

    selectRadioGroup(wrapper).setValue(KPI_PIE_CHART_SHOW_MODES.percent);

    expect(wrapper).toEmit('input', KPI_PIE_CHART_SHOW_MODES.percent);
  });

  test('Renders `pie-show-mode` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: KPI_PIE_CHART_SHOW_MODES.percent,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pie-show-mode` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: KPI_PIE_CHART_SHOW_MODES.numbers,
        name: 'custom_name',
        label: 'Custom label',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
