import Faker from 'faker';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { mockDateNow, mockModals } from '@unit/utils/mock-hooks';

import { MODALS } from '@/constants';

import PbehaviorRecurrenceRuleExceptionsField
  from '@/components/other/pbehavior/exceptions/fields/pbehavior-recurrence-rule-exceptions-field.vue';

const localVue = createVueInstance();

const stubs = {
  'pbehavior-exceptions-list': true,
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
const selectChooseExceptionButton = wrapper => selectButtonByIndex(wrapper, 1);

describe('pbehavior-recurrence-rule-exceptions-field', () => {
  const nowTimestamp = 1386435500000;
  mockDateNow(nowTimestamp);
  const $modals = mockModals();

  const factory = generateShallowRenderer(PbehaviorRecurrenceRuleExceptionsField, {
    localVue,
    stubs,
    mocks: { $modals },
  });
  const snapshotFactory = generateRenderer(PbehaviorRecurrenceRuleExceptionsField, { localVue, stubs });

  test('Exception added after trigger create button', () => {
    const exdates = [{ key: 'exdate-1', begin: 1, end: 2, type: '' }];
    const wrapper = factory({
      propsData: {
        exdates,
      },
    });

    selectAddExceptionButton(wrapper).vm.$emit('click');

    expect(wrapper).toEmit('input', [
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

    selectChooseExceptionButton(wrapper).vm.$emit('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.selectExceptionsLists,
        config: {
          exceptions,
          action: expect.any(Function),
        },
      },
    );

    const [{ config }] = $modals.show.mock.calls[0];
    const newExceptions = [
      ...exceptions,
      { name: Faker.datatype.string() },
    ];
    config.action(newExceptions);

    expect(wrapper).toEmit('update:exceptions', newExceptions);
  });

  test('Renders `pbehavior-exceptions-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
