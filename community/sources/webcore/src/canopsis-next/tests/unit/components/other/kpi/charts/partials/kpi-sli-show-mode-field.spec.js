import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { KPI_SLI_GRAPH_DATA_TYPE } from '@/constants';

import KpiSliShowModeField from '@/components/other/kpi/charts/partials/kpi-sli-show-mode-field';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(KpiSliShowModeField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiSliShowModeField, {
  localVue,

  ...options,
});

describe('kpi-sli-show-mode-field', () => {
  it('Unit changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: KPI_SLI_GRAPH_DATA_TYPE.time,
      },
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.setValue(KPI_SLI_GRAPH_DATA_TYPE.percent);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(KPI_SLI_GRAPH_DATA_TYPE.percent);
  });

  it('Renders `kpi-sli-show-mode-field` without props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: KPI_SLI_GRAPH_DATA_TYPE.time,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
