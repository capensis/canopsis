import { generateShallowRenderer } from '@unit/utils/vue';

import PbehaviorRecurrenceRuleExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

const pbehaviorExceptionsListStub = {
  name: 'PbehaviorExceptionsList',
  props: ['exceptions'],
  template: '<div />',
};
const pbehaviorExceptionsFieldStub = {
  name: 'PbehaviorExceptionsField',
  props: ['exdates'],
  template: '<div />',
};

const stubs = {
  'pbehavior-exceptions-list': pbehaviorExceptionsListStub,
  'pbehavior-exceptions-field': pbehaviorExceptionsFieldStub,
};

const selectPbehaviorExceptionsList = wrapper => wrapper.findComponent(pbehaviorExceptionsListStub);
const selectPbehaviorExceptionsField = wrapper =>
  wrapper.findComponent(pbehaviorExceptionsFieldStub);

describe('pbehavior-recurrence-rule-exceptions-field', () => {
  const factory = generateShallowRenderer(PbehaviorRecurrenceRuleExceptionsField, { stubs });

  test('Exceptions list shown when exceptions exist', () => {
    const exceptions = [{ name: 'exception-1' }];
    const wrapper = factory({
      propsData: {
        exceptions,
      },
    });

    expect(selectPbehaviorExceptionsList(wrapper).exists()).toBe(true);
    expect(selectPbehaviorExceptionsList(wrapper).props('exceptions')).toEqual(exceptions);
  });

  test('Exceptions list hidden when exceptions are empty', () => {
    const wrapper = factory({
      propsData: {
        exceptions: [],
      },
    });

    expect(selectPbehaviorExceptionsList(wrapper).exists()).toBe(false);
  });

  test('Passes exdates to exceptions field', () => {
    const exdates = [{ key: 'exdate-1' }];
    const wrapper = factory({
      propsData: {
        exdates,
      },
    });

    expect(selectPbehaviorExceptionsField(wrapper).props('exdates')).toEqual(exdates);
  });
});
