import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import { SERVICE_WEATHER_WIDGET_MODAL_TYPES } from '@/constants';

import FieldModalType from '@/components/sidebars/service-weather/form/fields/modal-type.vue';

const stubs = {
  'widget-settings-item': true,
  'v-radio-group': createInputStub('v-radio-group'),
};

const snapshotStubs = {
  'widget-settings-item': true,
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('field-modal-type', () => {
  const factory = generateShallowRenderer(FieldModalType, { stubs });

  const snapshotFactory = generateRenderer(FieldModalType, { stubs: snapshotStubs });

  test('Value changed after trigger radio group', () => {
    const wrapper = factory();

    selectRadioGroup(wrapper).triggerCustomEvent('input', SERVICE_WEATHER_WIDGET_MODAL_TYPES.both);

    expect(wrapper).toEmit('input', SERVICE_WEATHER_WIDGET_MODAL_TYPES.both);
  });

  test('Renders `field-modal-type` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `field-modal-type` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList,
        columns: [{
          value: 'column',
          label: 'Column',
        }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
