import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';

import PbehaviorRecurrenceRuleExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

const stubs = {
  'pbehavior-exceptions-list': true,
  'pbehavior-recurrence-rule-exceptions-list-menu': true,
  'pbehavior-exceptions-field': {
    template: `
    <div>
      <slot name="no-data" />
      <slot name="actions" />
    </div>
  `,
  },
  'c-alert': true,
};

const selectButtonByIndex = (wrapper, index) => wrapper.findAll('v-btn-stub').at(index);
const selectAddExceptionButton = wrapper => selectButtonByIndex(wrapper, 0);
const selectChooseExceptionButton = wrapper => wrapper.find('pbehavior-recurrence-rule-exceptions-list-menu-stub');

describe('pbehavior-recurrence-rule-exceptions-field', () => {
  const nowTimestamp = 1386435500000;
  mockDateNow(nowTimestamp);
  const $modals = mockModals();

  const factory = generateShallowRenderer(PbehaviorRecurrenceRuleExceptionsField, {
    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(PbehaviorRecurrenceRuleExceptionsField, { stubs });

  test('Exception added after trigger create button', () => {
    const exdates = [{ key: 'exdate-1', begin: 1, end: 2, type: '' }];
    const wrapper = factory({
      propsData: {
        exdates,
      },
    });

    selectAddExceptionButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmitInput([
      ...exdates,
      {
        key: expect.any(String),
        begin: new Date(1386370800000),
        end: new Date(1386370800000),
        type: '',
      },
    ]);
  });

  test('Exceptions selected after trigger select button', () => {
    const exceptions = [{
      name: Faker.datatype.string(),
    }];
    const wrapper = factory({
      propsData: {
        exceptions,
      },
    });
    const mewExceptions = [
      ...exceptions,
      {
        name: Faker.datatype.string(),
      },
    ];

    selectChooseExceptionButton(wrapper).triggerCustomEvent('input', mewExceptions);

    expect(wrapper).toEmit('update:exceptions', mewExceptions);
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
