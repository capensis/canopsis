import { createButtonStub } from '@unit/stubs/button';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

const stubs = {
  'v-btn': createButtonStub('v-btn'),
};

describe('chart-export-actions', () => {
  const factory = generateShallowRenderer(ChartExportActions, { stubs });

  it('Export csv event emitted', () => {
    const wrapper = factory({
      propsData: {
        chart: {},
      },
    });

    wrapper.findAll('button.v-btn').at(0).trigger('click');

    expect(wrapper).toHaveBeenEmit('export:csv');
  });

  it('Export png event emitted with correct blob', async () => {
    const canvas = document.createElement('canvas');
    const blob = new Blob();

    const wrapper = factory({
      propsData: {
        chart: {
          canvas,
        },
      },
    });

    const toBlobSpy = vi.spyOn(canvas, 'toBlob').mockImplementation(callback => callback(blob));

    const exportPngButton = wrapper.findAll('button.v-btn').at(1);

    await exportPngButton.trigger('click');
    await flushPromises();

    expect(wrapper).toEmit('export:png', blob);

    toBlobSpy.mockRestore();
  });
});
