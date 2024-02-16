import { generateRenderer } from '@unit/utils/vue';

import CInformationBlock from '@/components/common/block/c-information-block.vue';
import KpiCharts from '@/components/other/kpi/charts/kpi-charts';
import CHelpIcon from '@/components/common/icons/c-help-icon.vue';

const snapshotStubs = {
  'c-information-block': CInformationBlock,
  'c-help-icon': CHelpIcon,
  'kpi-alarms': true,
  'kpi-rating': true,
  'kpi-sli': true,
};

describe('kpi-charts', () => {
  const snapshotFactory = generateRenderer(KpiCharts, {
    stubs: snapshotStubs,
    attachTo: document.body,
  });

  it('Renders `kpi-charts` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });
});
