import moment from 'moment';

import { ACTION_TYPES, ACTION_AUTHOR } from '@/constants';

import { unsetInSeveralWithConditions } from '@/helpers/immutable';
import formToPbehavior from '@/helpers/forms/pbehavior/form-to-pbehavior';
import commentsToPbehaviorComments from '@/helpers/forms/pbehavior/comments-to-pbehavior-comments';
import exdatesToPbehaviorExdates from '@/helpers/forms/pbehavior/exdates-to-pbehavior-exdates';

export default function ({ generalParameters = {}, pbehaviorParameters = {}, snoozeParameters = {} }) {
  let data = { ...generalParameters };

  const patternsCondition = value => !value || !value.length;

  data = unsetInSeveralWithConditions(data, {
    'hook.event_patterns': patternsCondition,
    'hook.alarm_patterns': patternsCondition,
    'hook.entity_patterns': patternsCondition,
  });

  if (generalParameters.type === ACTION_TYPES.snooze) {
    const duration = moment.duration(
      parseInt(snoozeParameters.duration.duration, 10),
      snoozeParameters.duration.durationType,
    ).asSeconds();

    data.parameters = {
      ...snoozeParameters,
      duration,
    };
  } else if (generalParameters.type === ACTION_TYPES.pbehavior) {
    const pbehavior = formToPbehavior(pbehaviorParameters.general);

    pbehavior.comments =
      commentsToPbehaviorComments(pbehaviorParameters.comments);
    pbehavior.exdate = exdatesToPbehaviorExdates(pbehaviorParameters.exdate);

    data.parameters = { ...pbehavior };
  }

  data.parameters.author = ACTION_AUTHOR;

  return data;
}
