/*
We regroup this two function, because in the component, stepTitle can't access stepType since the 'this' doesn't exist
in the component filter
 */
import { ENTITY_INFOS_TYPE } from '@/constants';
import { capitalize } from 'lodash';

export function stepType(title) {
  if (title.startsWith('status')) {
    return ENTITY_INFOS_TYPE.status;
  } else if (title.startsWith('state')) {
    return ENTITY_INFOS_TYPE.state;
  }
  return ENTITY_INFOS_TYPE.action;
}

export function stepTitle(title, author) {
  let formattedStepTitle = '';
  if (stepType(title) !== ENTITY_INFOS_TYPE.action) {
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

