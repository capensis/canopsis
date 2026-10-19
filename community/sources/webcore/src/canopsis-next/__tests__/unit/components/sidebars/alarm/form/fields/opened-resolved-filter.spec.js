import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARMS_OPENED_VALUES } from '@/constants';

import OpenedResolvedFilter from '@/components/sidebars/alarm/form/fields/opened-resolved-filter.vue';

const stubs = {
  'widget-settings-item': true,
  'c-radio-group': true,
};

const selectRadioGroup = wrapper => wrapper.find('c-radio-group-stub');

describe('opened-resolved-filter', () => {
  const factory = generateShallowRenderer(OpenedResolvedFilter, { stubs });

  it.each(Object.entries(ALARMS_OPENED_VALUES))(
    'Value changed after to %s after trigger a radio field',
    (key, value) => {
      const wrapper = factory({
        propsData: {
          value:
            value === ALARMS_OPENED_VALUES.opened
              ? ALARMS_OPENED_VALUES.resolved
              : ALARMS_OPENED_VALUES.opened,
        },
      });

      selectRadioGroup(wrapper).triggerCustomEvent('input', String(value));

      expect(wrapper).toEmitInput(value);
    },
  );
});
