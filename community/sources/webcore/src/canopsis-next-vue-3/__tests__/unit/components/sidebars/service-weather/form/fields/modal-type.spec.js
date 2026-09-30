import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { SERVICE_WEATHER_WIDGET_MODAL_TYPES } from '@/constants';

import FieldModalType from '@/components/sidebars/service-weather/form/fields/modal-type.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('field-modal-type', () => {
  const factory = generateShallowRenderer(FieldModalType, { stubs });

  test('Value changed after trigger radio group', () => {
    const wrapper = factory();

    selectRadioGroup(wrapper).triggerCustomEvent('input', SERVICE_WEATHER_WIDGET_MODAL_TYPES.both);

    expect(wrapper).toEmitInput(SERVICE_WEATHER_WIDGET_MODAL_TYPES.both);
  });
});
