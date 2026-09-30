import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_FIELDS, ALARM_PATTERN_FIELDS, PATTERN_OPERATORS } from '@/constants';

import CAlarmPatternsField from '@/components/forms/fields/alarm/c-alarm-patterns-field.vue';

const stubs = {
  'pattern-editor-field': true,
};

describe('c-alarm-patterns-field', () => {
  const factory = generateShallowRenderer(CAlarmPatternsField, { stubs });

  test('Failed ticket field is represented by the ticket attribute operator', () => {
    const wrapper = factory({
      propsData: {
        patterns: {},
        attributes: [
          {
            value: ALARM_PATTERN_FIELDS.ticket,
            options: { disabled: false },
          },
          {
            value: ALARM_FIELDS.failedTicket,
            options: { disabled: false },
          },
        ],
      },
    });

    const ticketAttribute = wrapper.vm.availableAlarmAttributes.find(
      attribute => attribute.value === ALARM_PATTERN_FIELDS.ticket,
    );

    expect(ticketAttribute.options.operators).toContain(PATTERN_OPERATORS.ticketCreationFailed);
    expect(wrapper.vm.availableAlarmAttributes).not.toContainEqual(
      expect.objectContaining({ value: ALARM_FIELDS.failedTicket }),
    );
  });

  test('Multiple initiator values are displayed as closable chips', () => {
    const wrapper = factory({
      propsData: {
        patterns: {},
      },
    });

    const props = wrapper.vm.initiatorOptions.valueField.props({ value: ['user'] });

    expect(props).toEqual(
      expect.objectContaining({
        chips: true,
        multiple: true,
        closableChips: true,
      }),
    );
  });
});
