import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';

import KpiChartExportActions from '@/components/other/kpi/charts/partials/kpi-chart-export-actions';

const localVue = createVueInstance();

const stubs = {
  'v-btn': createButtonStub('v-btn'),
};

const factory = (options = {}) => shallowMount(KpiChartExportActions, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiChartExportActions, {
  localVue,

  ...options,
});

describe('kpi-chart-export-actions', () => {
  it('Export csv event emitted', () => {
    const wrapper = factory({
      propsData: {
        chart: {},
      },
    });

    const exportCsvButton = wrapper.findAll('button.v-btn').at(0);

    exportCsvButton.trigger('click');

    const exportCsvEvents = wrapper.emitted('export:csv');

    expect(exportCsvEvents).toHaveLength(1);
  });

  it('Export png event emitted with correct blob', (done) => {
    const canvas = document.createElement('canvas');

    const originalToBlob = canvas.toBlob;

    const wrapper = factory({
      propsData: {
        chart: {
          canvas,
        },
      },
    });

    const toBlobSpy = jest.spyOn(canvas, 'toBlob')
      .mockImplementation(callback => originalToBlob.call(canvas, (...args) => {
        callback(...args);

        localVue.nextTick(() => {
          const exportCsvEvents = wrapper.emitted('export:png');

          expect(exportCsvEvents).toHaveLength(1);

          const [exportEventData] = exportCsvEvents[0];

          expect(exportEventData).toEqual(expect.any(Blob));

          done();
        });
      }));

    const exportPngButton = wrapper.findAll('button.v-btn').at(1);

    exportPngButton.trigger('click');

    toBlobSpy.mockReset();
  });

  it('Renders `kpi-chart-export-actions` without props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        chart: {},
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `kpi-chart-export-actions` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        chart: {},
        downloading: true,
      },
    });

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});