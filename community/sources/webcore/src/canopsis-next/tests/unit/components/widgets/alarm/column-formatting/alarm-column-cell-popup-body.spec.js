import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import AlarmColumnCellPopupBody from '@/components/widgets/alarm/columns-formatting/alarm-column-cell-popup-body.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

const stubs = {
  'c-compiled-template': CCompiledTemplate,
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

    selectCloseButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmit('close');
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

    expect(wrapper).toMatchSnapshot();
  });
});
