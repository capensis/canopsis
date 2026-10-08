import { createButtonStub } from '@unit/stubs/button';
import { generateShallowRenderer } from '@unit/utils/vue';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import AlarmColumnCellPopupBody from '@/components/widgets/alarm/columns-formatting/alarm-column-cell-popup-body.vue';

const stubs = {
  'c-compiled-template': CCompiledTemplate,
  'c-runtime-template': CRuntimeTemplate,
  'v-card': { template: '<div><slot /></div>' },
  'v-card-title': { template: '<div><slot /></div>' },
  'v-btn': createButtonStub('v-btn'),
};

const selectCloseButton = wrapper => wrapper.find('button.v-btn');

describe('alarm-column-cell-popup-body', () => {
  const factory = generateShallowRenderer(AlarmColumnCellPopupBody, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Popup template closed after click on the button', async () => {
    const wrapper = factory({
      propsData: {
        template: '',
        alarm: {},
      },
    });

    await selectCloseButton(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('close');
  });
});
