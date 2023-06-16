import flushPromises from 'flush-promises';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import AlarmColumnCellPopupBody from '@/components/widgets/alarm/columns-formatting/alarm-column-cell-popup-body.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
};

const selectCloseButton = wrapper => wrapper.find('v-btn-stub');

describe('alarm-column-cell-popup-body', () => {
  const factory = generateShallowRenderer(AlarmColumnCellPopupBody, { stubs });
  const snapshotFactory = generateRenderer(AlarmColumnCellPopupBody, { stubs });

  it('Popup template closed after click on the button', async () => {
    const wrapper = factory({
      propsData: {
        template: '',
        alarm: {},
      },
    });

    const closeButton = selectCloseButton(wrapper);

    closeButton.vm.$emit('click');

    const closeEvents = wrapper.emitted('close');

    expect(closeEvents).toHaveLength(1);
  });

  it('Renders `alarm-column-cell-popup-body` with full alarm', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<p>Test</p><p>{{ alarm.name }}</p><p>{{ entity.name }}</p>',
        alarm: {
          name: 'alarm-name',
          entity: {
            name: 'entity-name',
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
