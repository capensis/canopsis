import { ALARM_STATES } from '@/constants';

import { formToChangeStateEvent } from '@/helpers/entities/alarm/form';

describe('alarm form helpers', () => {
  it('Converts change state form to event request', () => {
    const output = 'Change state note';

    expect(formToChangeStateEvent({
      state: ALARM_STATES.critical,
      output,
    })).toEqual({
      state: ALARM_STATES.critical,
      comment: output,
    });
  });

  it('Keeps existing change state comment when output is missing', () => {
    const comment = 'Existing comment';

    expect(formToChangeStateEvent({
      state: ALARM_STATES.major,
      comment,
    })).toEqual({
      state: ALARM_STATES.major,
      comment,
    });
  });
});
