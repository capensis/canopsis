import {
  CUSTOM_WIDGET_TEMPLATE,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import {
  alarmListAssociateTicketParametersToForm,
  formToAlarmListAssociateTicketParameters,
} from '@/helpers/entities/widget/forms/alarm';
import {
  formToWidgetTemplate,
  widgetTemplateToForm,
} from '@/helpers/entities/widget/template/form';
import { prepareWidgetTemplates } from '@/helpers/entities/widget/template/prepare';

describe('ticket format settings helpers', () => {
  test('uses custom form values until widget templates are loaded', () => {
    const form = alarmListAssociateTicketParametersToForm();

    expect(form.formatsTemplate).toBe(CUSTOM_WIDGET_TEMPLATE);
    expect(form.formats).toEqual([]);
  });

  test('preserves explicitly selected custom mode', () => {
    const form = alarmListAssociateTicketParametersToForm({
      formatsTemplate: '',
      formats: ['^INC-[0-9]+$'],
    });

    expect(form.formatsTemplate).toBe(CUSTOM_WIDGET_TEMPLATE);
    expect(form.formats.map(({ value }) => value)).toEqual(['^INC-[0-9]+$']);
  });

  test('uses custom mode when a deleted template left copied formats without a reference', () => {
    const form = alarmListAssociateTicketParametersToForm({ formats: ['^INC-[0-9]+$'] });

    expect(form.formatsTemplate).toBe(CUSTOM_WIDGET_TEMPLATE);
    expect(form.formats.map(({ value }) => value)).toEqual(['^INC-[0-9]+$']);
  });

  test('converts custom form values to the backend shape', () => {
    expect(
      formToAlarmListAssociateTicketParameters({
        formatsTemplate: CUSTOM_WIDGET_TEMPLATE,
        formats: [{ key: 'format-key', value: '^INC-[0-9]+$' }],
      }),
    ).toEqual({
      formatsTemplate: '',
      formats: ['^INC-[0-9]+$'],
    });
  });

  test('converts ticket format widget templates in both directions', () => {
    const form = widgetTemplateToForm({
      title: 'Jira ticket formats',
      type: WIDGET_TEMPLATES_TYPES.ticketFormats,
      formats: ['^INC-[0-9]+$', '^REQ-[0-9]+$'],
    });

    expect(form.formats.map(({ value }) => value)).toEqual([
      '^INC-[0-9]+$',
      '^REQ-[0-9]+$',
    ]);
    expect(formToWidgetTemplate(form)).toEqual({
      title: 'Jira ticket formats',
      type: WIDGET_TEMPLATES_TYPES.ticketFormats,
      formats: ['^INC-[0-9]+$', '^REQ-[0-9]+$'],
    });
  });

  test('prepares ticket format templates for selection in widget settings', () => {
    const [template] = prepareWidgetTemplates([
      {
        _id: 'ticket-formats-template-id',
        type: WIDGET_TEMPLATES_TYPES.ticketFormats,
        formats: ['^INC-[0-9]+$'],
      },
    ]);

    expect(template.formats.map(({ value }) => value)).toEqual(['^INC-[0-9]+$']);
  });
});
