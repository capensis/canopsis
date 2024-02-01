import { ENTITY_TYPES, ROOT_CAUSE_DIAGRAM_EVENTS_NODE_SIZE, ROOT_CAUSE_DIAGRAM_NODE_SIZE } from '@/constants';

import { getMapEntityText } from '@/helpers/entities/map/list';

// eslint-disable-next-line import/no-webpack-loader-syntax
import engineeringIcon from '!!svg-inline-loader?modules!@/assets/images/engineering.svg';

import { getEntityColor } from './color';

/**
 * Create vuetify icon element
 *
 * @param {string} name
 * @return {HTMLElement}
 */
const getIconElement = (name) => {
  const badgeIconEl = document.createElement('i');
  badgeIconEl.classList.add(
    'v-icon',
    'material-icons',
    'theme--light',
    'white--text',
  );

  badgeIconEl.innerHTML = name;

  badgeIconEl.style.width = '100%';
  badgeIconEl.style.height = '100%';
  badgeIconEl.style.borderRadius = '50%';

  return badgeIconEl;
};

/**
 * Create vuetify badge element
 *
 * @return {HTMLSpanElement}
 */
const getBadgeElement = () => {
  const badgeEl = document.createElement('span');
  badgeEl.classList.add(
    'v-badge__badge',
    'd-inline-flex',
    'justify-center',
    'align-center',
    'grey',
    'darken-1',
    'cursor-pointer',
    'pa-0',
  );
  badgeEl.style.width = '20px';
  badgeEl.style.height = '20px';

  return badgeEl;
};

/**
 * Create icon by entity fot state settings
 *
 * @return {HTMLSpanElement}
 */
export const getStateSettingsNodeIconElement = (node) => {
  const { isEvents, entity } = node;
  const size = `${isEvents ? 20 : 30}px`;

  const icon = isEvents
    ? 'textsms'
    : {
      [ENTITY_TYPES.service]: engineeringIcon,
      [ENTITY_TYPES.resource]: 'perm_identity',
      [ENTITY_TYPES.component]: 'developer_board',
    }[entity.type];

  const element = getIconElement(icon);

  element.style.fontSize = size;

  const svgElement = element.querySelector('svg');

  if (svgElement) {
    svgElement.style.width = size;
  }

  return element;
};

/**
 * Create vuetify progress element
 *
 * @return {HTMLDivElement}
 */
export const getProgressElement = () => {
  const progressContentCircleEl = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
  progressContentCircleEl.classList.add('v-progress-circular__overlay');
  progressContentCircleEl.setAttribute('fill', 'transparent');
  progressContentCircleEl.setAttribute('cx', '45.714285714285715');
  progressContentCircleEl.setAttribute('cy', '45.714285714285715');
  progressContentCircleEl.setAttribute('r', '15');
  progressContentCircleEl.setAttribute('stroke-width', '3');
  progressContentCircleEl.setAttribute('stroke-dasharray', '125.664');
  progressContentCircleEl.setAttribute('stroke-dashoffset', '125.66370614359172px');

  const progressContentEl = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  progressContentEl.setAttribute('viewBox', '22.857142857142858 22.857142857142858 45.714285714285715 45.714285714285715');
  progressContentEl.appendChild(progressContentCircleEl);

  const progressEl = document.createElement('div');
  progressEl.appendChild(progressContentEl);
  progressEl.classList.add(
    'v-progress-circular',
    'v-progress-circular--indeterminate',
    'v-progress-circular--visible',
    'white--text',
    'position-relative',
  );

  progressEl.style.width = '20px';
  progressEl.style.height = '20px';

  return progressEl;
};

/**
 * Create entity node button html
 *
 * @property {Object} node
 * @return {string}
 */
export const getEntityNodeElementHTML = (node) => {
  const { entity, pending, isEvents, opened } = node;

  const nodeSize = isEvents ? ROOT_CAUSE_DIAGRAM_EVENTS_NODE_SIZE : ROOT_CAUSE_DIAGRAM_NODE_SIZE;

  const nodeLabelEl = document.createElement('div');
  nodeLabelEl.classList.add('position-absolute');
  nodeLabelEl.style.top = `${nodeSize}px`;

  if (!isEvents) {
    nodeLabelEl.textContent = getMapEntityText(entity);
  }

  const nodeEl = document.createElement('div');
  nodeEl.appendChild(getStateSettingsNodeIconElement(node));
  nodeEl.appendChild(nodeLabelEl);
  nodeEl.classList.add('v-btn__content', 'position-relative', 'border-radius-rounded');
  nodeEl.style.width = `${nodeSize}px`;
  nodeEl.style.height = `${nodeSize}px`;
  nodeEl.style.justifyContent = 'center';
  nodeEl.style.background = getEntityColor(entity);

  if (pending || entity.depends_count > 0) {
    const badge = getBadgeElement();
    badge.dataset.id = entity._id;

    badge.appendChild(
      pending ? getProgressElement() : getIconElement(opened ? 'remove' : 'add'),
    );

    nodeEl.appendChild(badge);
  }

  return nodeEl.outerHTML;
};

/**
 * Create vuetify button html
 *
 * @property {Object} node
 * @return {string}
 */
export const getButtonHTML = (text) => {
  const btnContentEl = document.createElement('div');
  btnContentEl.classList.add('v-btn__content');
  btnContentEl.textContent = text;

  const btnEl = document.createElement('button');
  btnEl.classList.add(
    'v-btn',
    'v-btn--round',
    'theme--light',
  );
  btnEl.appendChild(btnContentEl);

  return btnEl.outerHTML;
};
