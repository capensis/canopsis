import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import PbehaviorRecurrenceRuleExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

const stubs = {
  'pbehavior-exceptions-list': true,
  'pbehavior-exceptions-field': true,
};

const selectPbehaviorExceptionsList = wrapper => wrapper.find('pbehavior-exceptions-list-stub');
const selectPbehaviorExceptionsField = wrapper => wrapper.find('pbehavior-exceptions-field-stub');

describe('pbehavior-recurrence-rule-exceptions-field', () => {
  const factory = generateShallowRenderer(PbehaviorRecurrenceRuleExceptionsField, { stubs });
  const snapshotFactory = generateRenderer(PbehaviorRecurrenceRuleExceptionsField, { stubs });

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

  test('Renders `pbehavior-exceptions-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehavior-exceptions-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        exdates: [
          { key: 'exdate-1' },
          { key: 'exdate-2' },
        ],
        exceptions: [
          { key: 'exception-1' },
        ],
        withExdateType: true,
        disabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
