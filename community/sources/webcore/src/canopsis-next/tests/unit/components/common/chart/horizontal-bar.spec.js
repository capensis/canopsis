import { mount, createVueInstance } from '@unit/utils/vue';

import HorizontalBar from '@/components/common/chart/horizontal-bar.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(HorizontalBar, {
  localVue,
  attachTo: document.body,

  ...options,
});

describe('horizontal-bar', () => {
  const labels = [1, 2, 3, 4];
  const data = [1, 2, 3, 4];
  const datasets = [{
    label: 'Dataset',
    backgroundColor: '#000',
    data,
  }];
  const options = {
    animation: false,
    responsive: false,
    plugins: {
      legend: {
        display: false,
      },
    },
  };
  const updatedLabels = [1, 2, 3, 4, 5, 6, 7, 8];
  const updatedData = [1, 2, 3, 4, 5, 6, 7, 8];
  const updatedDatasets = [{
    label: 'Updated dataset',
    backgroundColor: '#fff',
    data: updatedData,
  }];
  const updatedOptions = {
    animation: false,
    responsive: false,
    plugins: {
      legend: {
        position: 'right',
        labels: {
          boxWidth: 20,
        },
      },
    },
  };

  it('Renders `horizontal-bar` with default props and options.', () => {
    const wrapper = snapshotFactory({
      propsData: {
        options,
      },
    });

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `horizontal-bar` with custom props.', () => {
    const wrapper = snapshotFactory({
      propsData: {
        labels,
        datasets,
        options,
      },
    });

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `horizontal-bar` after update data.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        labels,
        datasets,
        options,
      },
    });

    await wrapper.setProps({
      labels: updatedLabels,
      datasets: updatedDatasets,
    });

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `horizontal-bar` after update options.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        labels,
        datasets,
        options,
      },
    });

    await wrapper.setProps({
      options: updatedOptions,
    });

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });
});
