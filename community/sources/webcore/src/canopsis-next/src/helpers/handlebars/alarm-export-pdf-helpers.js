import Handlebars from 'handlebars';

import { convertDateToString } from '../date/date';

import { harmonizeLinks } from '../links';

const CELL_BORDER = 'solid #f5f5f5 2px';

/**
 * Create paragraph HTML element
 *
 * @param {string} [content = '']
 * @returns {HTMLParagraphElement}
 */
export const createParagraph = (content = '') => {
  const p = document.createElement('p');

  p.style.padding = '12px';
  p.style.margin = '0px';
  p.innerHTML = content?.outerHTML ?? content;

  return p;
};

/**
 * Create table cell HTML element
 *
 * @param {string} [content = '']
 * @param {string} [lastRow]
 * @param {string} [lastCell]
 * @returns {HTMLTableCellElement}
 */
export const createTableCell = (content = '', lastRow, lastCell) => {
  const td = document.createElement('td');

  td.style.padding = '0px';
  td.style.width = '50%';
  td.style.backgroundColor = '#fff';

  if (!lastRow) {
    td.style.borderBottom = CELL_BORDER;
  }

  if (!lastCell) {
    td.style.borderRight = CELL_BORDER;
  }

  td.appendChild(createParagraph(content));

  return td;
};

/**
 * Create table row HTML element
 *
 * @param {string} [field = '']
 * @param {string} [value = '']
 * @param [last]
 * @returns {HTMLTableRowElement}
 */
export const createTableRow = (field = '', value = '', last) => {
  const tr = document.createElement('tr');

  tr.appendChild(createTableCell(field, last));
  tr.appendChild(createTableCell(value, last, true));

  return tr;
};

/**
 * Create table HTML element
 *
 * @returns {HTMLTableElement}
 */
export function createTable() {
  const table = document.createElement('table');

  table.style.width = '100%';
  table.style.borderCollapse = 'collapse';
  table.style.tableLayout = 'fixed';
  table.style.marginBottom = '2px';
  table.style.marginRight = '2px';

  return table;
}

/**
 * Convert infos field to html string
 *
 * @returns {string}
 */
export function infos() {
  const table = createTable();
  const tbody = document.createElement('tbody');

  Object.values(this.infos ?? {}).forEach((info, rootIndex, rootArray) => {
    Object.entries(info).forEach(([key, value], index, array) => {
      const lastRow = index === array.length - 1 && rootIndex === rootArray.length - 1;

      tbody.appendChild(createTableRow(key, value, lastRow));
    });
  });

  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}

/**
 * Convert pbehavior_info field to html string
 *
 * @returns {string}
 */
export function pbehaviorInfo() {
  const { pbehavior_info: info = {} } = this;

  if (!info) {
    return '';
  }

  const table = createTable();
  const tbody = document.createElement('tbody');

  tbody.appendChild(createTableRow('Enter time', info.timestamp));
  tbody.appendChild(createTableRow('Name', info.name));
  tbody.appendChild(createTableRow('Type', info.type));
  tbody.appendChild(createTableRow('Reason', info.reason, true));
  table.appendChild(tbody);

  return new Handlebars.SafeString(table.outerHTML);
}
/**
 * Convert ticket_info field to html string
 *
 * @returns {string}
 */
export function ticketInfo() {
  return convertDateToString(this.current_date);
}
/**
 * Convert last_comment field to html string
 *
 * @returns {string}
 */
export function lastComment() {
  if (!this.last_comment) {
    return '';
  }

  const div = document.createElement('div');

  div.appendChild(createParagraph(this.last_comment.t));
  div.appendChild(createParagraph(this.last_comment.m));

  return new Handlebars.SafeString(div.outerHTML);
}

/**
 * Convert tags field to html string
 *
 * @returns {string}
 */
export function tags() {
  const div = document.createElement('div');

  (this.tags ?? []).forEach(tag => div.appendChild(createParagraph(tag)));

  return new Handlebars.SafeString(div.outerHTML);
}

/**
 * Convert links field to html string
 *
 * @returns {string}
 */
export function links() {
  const div = document.createElement('div');

  harmonizeLinks(this.links).forEach(link => div.appendChild(createParagraph(link.url)));

  return new Handlebars.SafeString(div.outerHTML);
}

export const createInstanceWithHelpers = () => {
  const instance = Handlebars.create();

  instance.registerHelper('infos', infos);
  instance.registerHelper('tags', tags);
  instance.registerHelper('links', links);
  instance.registerHelper('ticket_info', ticketInfo);
  instance.registerHelper('last_comment', lastComment);
  instance.registerHelper('pbehavior_info', pbehaviorInfo);

  return instance;
};
