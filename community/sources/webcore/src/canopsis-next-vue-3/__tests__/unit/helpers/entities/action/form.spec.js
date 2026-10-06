import { ACTION_TYPES } from '@/constants';

import { actionToForm, formToAction } from '@/helpers/entities/action/form';

describe('entity pause scenario action form', () => {
  test('converts pause parameters with author fields', () => {
    const form = actionToForm({
      type: ACTION_TYPES.pause,
      parameters: {
        output: 'Pause {{ .Entity.ID }}',
        type: 'maintenance',
        inherited: true,
        forward_author: false,
        author: 'Pause author',
      },
    });

    expect(form.parameters[ACTION_TYPES.pause]).toEqual({
      output: 'Pause {{ .Entity.ID }}',
      type: 'maintenance',
      inherited: true,
      forward_author: false,
      author: 'Pause author',
    });
    expect(formToAction(form).parameters).toEqual({
      output: 'Pause {{ .Entity.ID }}',
      type: 'maintenance',
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
