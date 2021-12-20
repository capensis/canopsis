import { mount, createVueInstance } from '@unit/utils/vue';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import KpiCharts from '@/components/other/kpi/charts/kpi-charts';
import CHelpIcon from '@/components/common/icons/c-help-icon.vue';

const localVue = createVueInstance();

const snapshotStubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': CHelpIcon,
  'kpi-alarms': true,
  'kpi-rating': true,
  'kpi-sli': true,
};

const snapshotFactory = (options = {}) => mount(KpiCharts, {
  localVue,
  stubs: snapshotStubs,
  attachTo: document.body,

  ...options,
});

describe('kpi-charts', () => {
  it('Renders `kpi-charts` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });
});
