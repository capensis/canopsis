import { generateShallowRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES } from '@/constants';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import AlarmColumnValue from '@/components/widgets/alarm/columns-formatting/alarm-column-value.vue';

const stubs = {
  'color-indicator-wrapper': true,
  'alarm-column-cell': true,
  'c-compiled-template': CCompiledTemplate,
  'c-runtime-template': CRuntimeTemplate,
};

const selectAlarmColumnCell = wrapper => wrapper.find('alarm-column-cell-stub');

describe('alarm-column-value', () => {
  const factory = generateShallowRenderer(AlarmColumnValue, { stubs });

  it('Click state emitted after trigger click state event on alarm cell', async () => {
    const wrapper = factory({
      propsData: {
        alarm: {
          entity: {},
        },
        widget: {},
        column: {
          colorIndicator: COLOR_INDICATOR_TYPES.state,
        },
      },
    });

    selectAlarmColumnCell(wrapper).triggerCustomEvent('click:state');

    expect(wrapper).toHaveBeenEmit('click:state');
  });
});
