import Vue from 'vue';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createButtonStub } from '@unit/stubs/button';

import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

const stubs = {
  'v-btn': createButtonStub('v-btn'),
};

describe('chart-export-actions', () => {
  const factory = generateShallowRenderer(ChartExportActions, { stubs });
  const snapshotFactory = generateRenderer(ChartExportActions);

  it('Export csv event emitted', () => {
    const wrapper = factory({
      propsData: {
        chart: {},
      },
    });

    wrapper.findAll('button.v-btn').at(0).trigger('click');

    expect(wrapper).toHaveBeenEmit('export:csv');
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

        Vue.nextTick(() => {
          expect(wrapper).toEmit('export:png', expect.any(Blob));

          done();
        });
      }));

    const exportPngButton = wrapper.findAll('button.v-btn').at(1);

    exportPngButton.trigger('click');

    toBlobSpy.mockReset();
  });

  it('Renders `chart-export-actions` without props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        chart: {},
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `chart-export-actions` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        chart: {},
        downloading: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
