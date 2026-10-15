import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { MODALS } from '@/constants';

import PbehaviorRecurrenceRuleField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-recurrence-rule-field.vue';

const stubs = {
  'c-label': true,
  'c-action-btn': true,
};

const selectCreateRruleButton = wrapper => wrapper.find('v-btn-stub');
const selectRemoveRruleButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);

describe('pbehavior-recurrence-rule-field', () => {
  const $modals = mockModals();

  const factory = generateShallowRenderer(PbehaviorRecurrenceRuleField, {
    stubs,
    mocks: { $modals },
  });

  test('Rrule created after trigger create button', () => {
    const rrule = '';
    const exdates = [];
    const exceptions = [];
    const wrapper = factory({
      propsData: {
        form: {
          rrule,
          exceptions,
          exdates,
        },
      },
    });

    selectCreateRruleButton(wrapper).triggerCustomEvent('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createRecurrenceRule,
      config: {
        rrule,
        exceptions,
        exdates,
        withExdateType: false,
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    const newRrule = 'FREQ=DAILY';
    const newExceptions = [{}];
    const newExdates = [{}, {}];

    config.action({
      rrule: newRrule,
      exceptions: newExceptions,
      exdates: newExdates,
    });

    expect(wrapper).toEmitInput({
      rrule: newRrule,
      exceptions: newExceptions,
      exdates: newExdates,
    });
  });

  test('Rrule removed after trigger remove button', () => {
    const wrapper = factory({
      propsData: {
        form: {
          rrule: 'FREQ=DAILY',
        },
      },
    });

    selectRemoveRruleButton(wrapper).triggerCustomEvent('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [{ config }] = $modals.show.mock.calls[0];

    config.action();

    expect(wrapper).toEmitInput({
      rrule: '',
    });
  });
});
