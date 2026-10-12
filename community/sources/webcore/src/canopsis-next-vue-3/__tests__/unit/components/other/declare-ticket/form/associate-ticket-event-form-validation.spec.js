import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import {
  computed, defineComponent, toRef, watch,
} from 'vue';

import { eventToAssociateTicketForm } from '@/helpers/entities/associate-ticket/event/form';

import { useField, useForm } from '@/hooks/form';

import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

const ValidatedTicketIdField = defineComponent({
  props: {
    value: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'ticket_id',
    },
    required: {
      type: Boolean,
      default: false,
    },
    validationRules: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:value'],

  setup(props, { emit }) {
    const validate = async (value) => {
      if (props.required && !value) {
        return 'required';
      }

      for (const rule of props.validationRules) {
        const result = await rule(value);

        if (result !== true) {
          return result;
        }
      }

      return true;
    };
    const { value, errorMessage, setValue } = useField(toRef(props, 'name'), validate, {
      initialValue: props.value,
      validateOnValueUpdate: false,
    });
    const inputValue = computed({
      get: () => value.value,
      set: (newValue) => {
        setValue(newValue, false);
        emit('update:value', newValue);
      },
    });

    watch(() => props.value, newValue => setValue(newValue, false));

    return {
      errorMessage,
      inputValue,
    };
  },

  template: `
    <div>
      <input v-model="inputValue" />
      <span>{{ errorMessage }}</span>
    </div>
  `,
});

const TestForm = defineComponent({
  components: { AssociateTicketEventForm },

  props: {
    ticketFormats: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['submit'],

  setup(props, { emit }) {
    const form = useForm({ initialValues: eventToAssociateTicketForm() });
    const submit = form.handleSubmit(() => emit('submit'));
    const updateForm = (value) => {
      form.values.value = value;
    };

    return {
      values: form.values,
      submit,
      updateForm,
    };
  },

  template: `
    <form @submit.prevent="submit">
      <associate-ticket-event-form
        :form="values"
        :ticket-formats="ticketFormats"
        @update:form="updateForm"
      />
    </form>
  `,
});

const stubs = {
  'c-form-block-row': {
    template: '<div><slot /></div>',
  },
  'c-name-field': true,
  'c-help-icon': true,
  'declare-ticket-rule-ticket-id-text-field': ValidatedTicketIdField,
  'declare-ticket-rule-ticket-url-text-field': true,
  'declare-ticket-rule-ticket-url-title-field': true,
  'declare-ticket-rule-ticket-custom-fields-field': true,
};

describe('associate-ticket-event-form validation', () => {
  const factory = generateRenderer(TestForm, { stubs });

  test('prevents submission until the ticket matches a widget format', async () => {
    const wrapper = factory({
      propsData: {
        ticketFormats: ['^INC-[0-9]+$'],
      },
    });
    const ticketInput = wrapper.find('input');

    await ticketInput.setValue('REQ-42');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(wrapper.emitted('submit')).toBeUndefined();
    expect(wrapper.text()).toContain('Incorrect ticket format');

    await ticketInput.setValue('INC-42');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(wrapper.emitted('submit')).toHaveLength(1);
  });
});
