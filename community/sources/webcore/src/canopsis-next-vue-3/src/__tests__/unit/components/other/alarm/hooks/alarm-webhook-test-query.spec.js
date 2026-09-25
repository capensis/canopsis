import { mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import { defineComponent, toRef } from 'vue';

import { declareTicketRuleErrorsToForm } from '@/helpers/entities/declare-ticket/rule/form';

import { useAlarmWebhookTestQuery } from '@/components/other/alarm/hooks/alarm-webhook-test-query';

const TestComponent = defineComponent({
  props: {
    form: {
      type: Object,
      required: true,
    },
    createExecution: {
      type: Function,
      required: true,
    },
  },
  setup(props) {
    return useAlarmWebhookTestQuery({
      form: toRef(props, 'form'),
      socketRoom: 'test-room',
      createExecution: props.createExecution,
      fetchExecution: vi.fn(),
      getExecutionData: ({ alarm, form }) => ({ alarm, ...form }),
      errorsToValidation: errors => declareTicketRuleErrorsToForm(errors, props.form),
    });
  },
  template: '<div />',
});

describe('useAlarmWebhookTestQuery', () => {
  const $popups = mockPopups();
  const factory = generateShallowRenderer(TestComponent, {
    mocks: { $popups },
  });

  it('maps backend field errors when test execution creation fails', async () => {
    const error = {
      'webhooks.0.declare_ticket.ticket_id': 'The field is a required field',
    };
    const fieldPath = 'webhooks.webhook-key.declare_ticket.ticket_id';
    const createExecution = vi.fn().mockRejectedValueOnce(error);
    const wrapper = factory({
      propsData: {
        form: {
          webhooks: [{ key: 'webhook-key', declare_ticket: { ticket_id: '' } }],
        },
        createExecution,
      },
    });
    globalThis.__TEST_FORM__.getAllPathStatesWithInstances.mockReturnValue([{ path: fieldPath }]);

    wrapper.vm.alarm = 'alarm-id';
    await wrapper.vm.runTestExecution();
    await flushPromises();

    expect(globalThis.__TEST_FORM__.addError).toHaveBeenCalledWith(
      fieldPath,
      'The field is a required field',
    );
    expect($popups.error).not.toHaveBeenCalled();
  });
});
