import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { KPI_SLI_GRAPH_DATA_TYPE } from '@/constants';

import KpiSliShowModeField from '@/components/other/kpi/charts/form/fields/kpi-sli-show-mode-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

describe('kpi-sli-show-mode-field', () => {
  const factory = generateShallowRenderer(KpiSliShowModeField, { stubs });

  it('Unit changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: KPI_SLI_GRAPH_DATA_TYPE.time,
      },
    });

    wrapper.find('select.v-select').setValue(KPI_SLI_GRAPH_DATA_TYPE.percent);

    expect(wrapper).toEmitInput(KPI_SLI_GRAPH_DATA_TYPE.percent);
  });
});
