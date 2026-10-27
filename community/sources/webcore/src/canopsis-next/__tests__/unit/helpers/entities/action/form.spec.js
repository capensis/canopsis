import {
  ACTION_TYPES,
  MIN_DURATION_VALUE,
  TIME_UNITS,
  WORKFLOW_TYPES,
} from '@/constants';

import {
  actionToForm,
  formToAction,
  snoozeActionParametersToForm,
} from '@/helpers/entities/action/form';

describe('scenario action form', () => {
  test('uses one second as the default snooze duration', () => {
    expect(snoozeActionParametersToForm().duration).toEqual({
      value: MIN_DURATION_VALUE,
      unit: TIME_UNITS.second,
    });
  });

  test('continues with remaining steps after a new webhook succeeds', () => {
    const webhookParameters = actionToForm().parameters[ACTION_TYPES.webhook];

    expect(webhookParameters.stop_on_success).toBe(WORKFLOW_TYPES.continue);
    expect(webhookParameters.stop_on_fail).toBe(WORKFLOW_TYPES.stop);
  });

  test('converts pause parameters with author fields', () => {
    const form = actionToForm({
      type: ACTION_TYPES.pause,
      parameters: {
        output: 'Pause {{ .Entity.ID }}',
        pause_type: 'maintenance',
        inherited: true,
        forward_author: false,
        author: 'Pause author',
      },
    });

    expect(form.parameters[ACTION_TYPES.pause]).toEqual({
      output: 'Pause {{ .Entity.ID }}',
      pause_type: 'maintenance',
      inherited: true,
      forward_author: false,
      author: 'Pause author',
    });
    expect(formToAction(form).parameters).toEqual({
      output: 'Pause {{ .Entity.ID }}',
      pause_type: 'maintenance',
      inherited: true,
      forward_author: false,
      author: 'Pause author',
    });
  });

  test('converts resume parameters with author fields', () => {
    const form = actionToForm({
      type: ACTION_TYPES.resume,
      parameters: {
        output: 'Resume {{ .Entity.ID }}',
        forward_author: false,
        author: 'Resume author',
      },
    });

    expect(form.parameters[ACTION_TYPES.resume]).toEqual({
      output: 'Resume {{ .Entity.ID }}',
      forward_author: false,
      author: 'Resume author',
    });
    expect(formToAction(form).parameters).toEqual({
      output: 'Resume {{ .Entity.ID }}',
      forward_author: false,
      author: 'Resume author',
    });
  });
});
