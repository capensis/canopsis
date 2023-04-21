import { isObject } from 'lodash';
import Handlebars from 'handlebars';

import { convertDateToString } from '@/helpers/date/date';

import { registerHelper, unregisterHelper } from './index';

const cellBorder = 'solid #f5f5f5 2px';

export const createParagraph = (content = '') => {
  const p = document.createElement('p');

  p.style.padding = '12px';
  p.style.margin = '0px';
  p.innerText = content;

  return p;
};

export const createLink = (url = '') => {
  const a = document.createElement('a');

  a.href = url;

  return a;
};

export const createTableValueCell = (content = '', { lastRow = false, lastCell = false } = {}) => {
  const td = document.createElement('td');
  td.style.padding = '0px';
  td.style.backgroundColor = '#fff';

  if (!lastRow) {
    td.style.borderBottom = cellBorder;
  }

  if (!lastCell) {
    td.style.borderRight = cellBorder;
  }

  td.appendChild(createParagraph(content));

  return td;
};

export const createValueTableBody = (object) => {
  const tbody = document.createElement('tbody');

  Object.entries(object).forEach(([key, value], index, array) => {
    const tr = document.createElement('tr');
    const lastRow = index === array.length - 1;

    const keyTd = createTableValueCell({ lastRow });
    const valueTd = createTableValueCell({ lastRow, lastColumn: true });

    keyTd.appendChild(createParagraph(key));
    valueTd.appendChild(isObject(value) ? createValueTableBody(value) : createParagraph(value));
    tr.appendChild(keyTd);
    tr.appendChild(valueTd);
    tbody.appendChild(tr);
  });

  return tbody;
};

export function createDeepTable() {
  const table = document.createElement('table');

  table.style.borderCollapse = 'collapse';

  return table;
}

export const createTableHead = (content) => {
  const thead = document.createElement('thead');
  const tr = document.createElement('tr');
  const th = document.createElement('th');

  th.colSpan = 2;
  th.style.fontSize = '18px';
  th.style.color = '#fff';
  th.style.backgroundColor = '#2fab63';
  th.style.border = 'solid #fff 2px';

  th.appendChild(createParagraph(content));
  tr.appendChild(th);
  thead.appendChild(tr);

  return thead;
};

export const createTableKeyCell = (content = '') => {
  const td = document.createElement('td');
  td.style.fontWeight = '700';
  td.style.padding = '0px !important';
  td.style.border = 'solid #fff 2px';
  td.style.borderTop = 'none';
  td.style.backgroundColor = '#f5f5f5';
  td.style.color = '#666';

  td.appendChild(createParagraph(content));

  return td;
};

export const createTableBody = (alarm) => {
  const tbody = document.createElement('tbody');

  Object.entries(alarm).forEach(([key, value]) => {
    const tr = document.createElement('tr');
    tr.appendChild(createTableKeyCell(key));
    tr.appendChild(createTableValueCell(value));

    tbody.appendChild(tr);
  });

  return tbody;
};

export const createTable = (obj) => {
  const div = document.createElement('div');
  const table = document.createElement('table');

  div.style.width = '1000px';

  table.style.width = '100%';
  table.style.fontFamily = 'Roboto, sans-serif';
  table.style.fontSize = '15px';
  table.style.lineHeight = '1.2';
  table.style.color = '#434343';
  table.style.overflowWrap = 'break-word';
  table.style.borderCollapse = 'collapse';

  table.appendChild(createTableHead('Alarm display'));
  table.appendChild(createTableBody(obj));

  div.appendChild(table);

  return div;
};

/**
 * Convert infos to html
 *
 * @returns {string}
 */
export function infos() {
  const table = createDeepTable();
  const tbody = document.createElement('tbody');

  Object.values(this.infos).forEach((info, rootIndex, rootArray) => {
    Object.entries(info).forEach(([key, value], index, array) => {
      const tr = document.createElement('tr');
      const lastRow = index === array.length - 1 && rootIndex === rootArray.length - 1;

      tr.appendChild(createTableValueCell(key, { lastRow }));
      tr.appendChild(createTableValueCell(value, { lastRow, lastCell: true }));
      tbody.appendChild(tr);
    });
  });

  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}

/**
 * Convert current_date to string
 *
 * @returns {string}
 */
export function pbehaviorInfo() { // TODO
  return convertDateToString(this.current_date);
}
/**
 * Convert current_date to string
 *
 * @returns {string}
 */
export function ticketInfo() { // TODO
  return convertDateToString(this.current_date);
}
/**
 * Convert comments to html
 *
 * @returns {string}
 */
export function lastComment() {
  const table = createDeepTable();
  const tbody = document.createElement('tbody');

  (this.comments ?? []).forEach((comment = {}) => {
    const tr = document.createElement('tr');
    const td = document.createElement('td');
    td.appendChild(createParagraph(convertDateToString(comment.created)));
    td.appendChild(createParagraph(comment.message));
    tr.appendChild(td);
    tbody.appendChild(tr);
  });

  return new Handlebars.SafeString(table.outerHTML);
}

/**
 * Convert current_date to string
 *
 * @returns {string}
 */
export function tags() { // TODO
  return convertDateToString(this.current_date);
}
/**
 * Convert current_date to string
 *
 * @returns {string}
 */
export function links() { // TODO
  return convertDateToString(this.current_date);
}

export const createRegistererAllAlarmHelpers = () => {
  const existsHelpers = {};
  const helpersForRegister = {
    infos,
    tags,
    links,
    last_comment: lastComment,
    pbehavior_info: pbehaviorInfo,
  };

  Object.entries(helpersForRegister).forEach(([name, helper]) => {
    if (Handlebars.helpers[name]) {
      existsHelpers[name] = Handlebars.helpers[name];
      unregisterHelper(name);
    }

    registerHelper(name, helper);
  });

  return () => {
    Object.keys(helpersForRegister).forEach(name => unregisterHelper(name));
    Object.entries(existsHelpers).forEach(([name, helper]) => registerHelper(name, helper));
  };
};
