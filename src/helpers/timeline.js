/*
We regroup this two function, because in the component, stepTitle can't access stepType since the 'this' doesn't exist
in the component filter
 */
import { STEPS_TYPES } from '@/constants';
import capitalize from 'lodash/capitalize';

export function stepType(title) {
  if (title.startsWith('status')) {
    return STEPS_TYPES.status;
  } else if (title.startsWith('state')) {
    return STEPS_TYPES.state;
  }
  return STEPS_TYPES.action;
}
export function stepTitle(title, author) {
  let formattedStepTitle = '';
  if (stepType(title) !== STEPS_TYPES.action) {
    formattedStepTitle = title.replace(/(status)|(state)/g, '$& ');
    formattedStepTitle = formattedStepTitle.replace(/(inc)|(dec)/g, '$&reased');
  } else {
    formattedStepTitle = title.replace(/(declare)|(ack)/g, '$& ');
  }
  formattedStepTitle += ' by ';
  if (author === 'canopsis.engine') {
    formattedStepTitle += 'system';
  } else {
    formattedStepTitle += author;
  }
  return capitalize(formattedStepTitle);
}

